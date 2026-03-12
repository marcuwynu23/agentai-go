package core

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"agentai-go/internal/config"
)

// Planner converts user goals into structured plans using AI.
type Planner struct {
	config  *config.Config
	aiCore  *AICore
	analyzer *CodebaseAnalyzer
}

// NewPlanner creates a new Planner.
func NewPlanner(cfg *config.Config) *Planner {
	return &Planner{
		config:  cfg,
		aiCore:  NewAICore(cfg),
		analyzer: NewCodebaseAnalyzer(""),
	}
}

// CreatePlan builds a plan from goal, memory, conversation context, and optional codebase analysis.
func (p *Planner) CreatePlan(ctx context.Context, goal string, memory map[string]interface{}, conversationContext, workspacePath string) (*Plan, error) {
	var codebaseAnalysis string
	if workspacePath != "" {
		p.analyzer.WorkspacePath = workspacePath
		analysis, err := p.analyzer.Analyze()
		if err == nil {
			codebaseAnalysis = p.analyzer.FormatAnalysisForAI(analysis)
		}
	}
	prompt := p.buildPlanPrompt(goal, memory, conversationContext, codebaseAnalysis)
	planResponse, err := p.aiCore.GeneratePlan(ctx, prompt)
	if err != nil {
		return nil, err
	}
	return p.parsePlanResponse(planResponse, goal)
}

func (p *Planner) buildPlanPrompt(goal string, memory map[string]interface{}, conversationContext, codebaseAnalysis string) string {
	memoryContext := "No previous actions"
	if hist, ok := memory["history"].([]interface{}); ok && len(hist) > 0 {
		last := hist
		if len(last) > 5 {
			last = last[len(last)-5:]
		}
		b, _ := json.MarshalIndent(last, "", "  ")
		memoryContext = "Previous actions: " + string(b)
	}
	conversationText := ""
	if conversationContext != "" && conversationContext != "No previous conversation." {
		conversationText = "\n\nPrevious Conversation:\n" + conversationContext
	}
	analysisText := "\n\nNote: This is a new project with no existing codebase."
	if codebaseAnalysis != "" {
		analysisText = "\n\n" + codebaseAnalysis + `

CRITICAL INSTRUCTIONS FOR EXISTING CODEBASE:
1. FIRST PRIORITY: Fix all critical bugs and security issues found in the analysis BEFORE adding new features
2. Review existing code structure, patterns, and conventions
3. Understand current dependencies and frameworks in use
4. Plan how new code integrates with existing architecture
5. Decide whether to modify existing files or create new ones
6. Ensure new code follows the same coding style and patterns
7. If bugs are detected, create code_generation steps to fix them as separate steps
8. Consider refactoring opportunities if code quality issues are found`
	}

	return `You are an expert code generation planner. Break down the following user goal into detailed, actionable steps.

User Goal: ` + goal + `

Project Context:
` + memoryContext + conversationText + analysisText + `

Generate a comprehensive plan with the following step types:
- file_creation: Create new files (use for initial file setup)
- code_generation: Generate or modify code (use for implementing functionality)
- test_creation: Create test cases (use after code generation)
- command_execution: Execute shell commands like npm install, npm init, etc.

IMPORTANT GUIDELINES:
1. Break down the goal into logical, sequential steps
2. Each step should be specific and actionable
3. Use file_creation for package.json, README, config files
4. Use code_generation for implementing actual functionality
5. Include test_creation steps after code is generated
6. Include command_execution for setup commands (npm install, etc.)
7. Set proper dependencies between steps
8. Use realistic file paths (relative to workspace root, no "workspace/" prefix)
9. If codebase analysis is provided, fix critical issues first, then integrate new code
10. Always consider existing dependencies and frameworks when planning
11. Maintain consistency with existing code patterns

For each step, provide:
- id: unique identifier (step_1, step_2, etc.)
- type: one of: file_creation, code_generation, test_creation, command_execution
- description: clear description of what this step accomplishes
- target: file path (e.g., "app.js", "package.json") or command (e.g., "npm init -y")
- dependencies: array of step IDs this depends on (empty array if no dependencies)

Return ONLY valid JSON (no markdown, no code blocks) with this structure:
{
  "goal": "` + goal + `",
  "steps": [
    {
      "id": "step_1",
      "type": "file_creation",
      "description": "Create package.json file",
      "target": "package.json",
      "dependencies": []
    }
  ]
}`
}

func (p *Planner) parsePlanResponse(response, goal string) (*Plan, error) {
	re := regexp.MustCompile(`\{[\s\S]*\}`)
	m := re.FindString(response)
	if m == "" {
		return fallbackPlan(goal), nil
	}
	var plan Plan
	if err := json.Unmarshal([]byte(m), &plan); err != nil {
		return fallbackPlan(goal), nil
	}
	if plan.Steps == nil {
		plan.Steps = []*PlanStep{}
	}
	for i, step := range plan.Steps {
		if step.ID == "" {
			step.ID = "step_" + fmt.Sprintf("%d", i+1)
		}
		if step.Type == "" {
			step.Type = "code_generation"
		}
		if step.Description == "" {
			step.Description = "Unnamed step"
		}
		if step.Dependencies == nil {
			step.Dependencies = []string{}
		}
	}
	return &plan, nil
}

func fallbackPlan(goal string) *Plan {
	return &Plan{
		Goal: goal,
		Steps: []*PlanStep{
			{
				ID:           "step_1",
				Type:         "code_generation",
				Description: "Generate code based on user goal",
				Target:       "",
				Dependencies: []string{},
			},
		},
	}
}
