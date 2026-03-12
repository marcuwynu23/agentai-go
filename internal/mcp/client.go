package mcp

import (
	"context"
	"regexp"
	"strings"

	"agentai-go/internal/types"
)

// CodeGenFunc generates code from a prompt (injected to avoid importing core).
type CodeGenFunc func(ctx context.Context, prompt string) (string, error)

// Client coordinates file, command, and test servers and uses AI for code generation.
type Client struct {
	WorkspacePath string
	FileServer    *FileServer
	CommandServer *CommandServer
	TestServer    *TestServer
	GenerateCode  CodeGenFunc
}

// NewClient creates an MCP client with in-process servers.
func NewClient(workspacePath string, genCode CodeGenFunc) *Client {
	return &Client{
		WorkspacePath: workspacePath,
		FileServer:    NewFileServer(workspacePath),
		CommandServer: NewCommandServer(workspacePath),
		TestServer:    NewTestServer(workspacePath),
		GenerateCode:  genCode,
	}
}

// StepResult is a generic step execution result.
type StepResult struct {
	Success  bool
	Error    string
	FilePath string
	Message  string
}

// HandleFileCreation creates a file, generating content via AI if needed.
func (c *Client) HandleFileCreation(ctx context.Context, step *types.Step, reasoning *types.Reasoning) StepResult {
	filePath := step.Target
	if filePath == "" {
		filePath = strings.ToLower(strings.ReplaceAll(step.Description, " ", "_")) + ".js"
	}
	content := step.Content
	if content == "" && reasoning != nil && reasoning.Instructions != "" {
		prompt := "Create a " + filePath + " file. " + reasoning.Instructions + "\n\nGenerate the complete file content:"
		gen, err := c.GenerateCode(ctx, prompt)
		if err == nil {
			content = cleanCodeBlock(gen)
		} else {
			content = "// " + step.Description + "\n// " + reasoning.Instructions + "\n"
		}
	}
	if content == "" {
		content = "// " + step.Description + "\n"
	}
	res := c.FileServer.HandleCreateFile(filePath, content)
	return StepResult{Success: res.Success, Error: res.Error, FilePath: res.FilePath, Message: res.Message}
}

// HandleCodeGeneration generates or updates code in target file.
func (c *Client) HandleCodeGeneration(ctx context.Context, step *types.Step, reasoning *types.Reasoning) StepResult {
	filePath := step.Target
	if filePath == "" {
		return StepResult{Success: false, Error: "no target file specified for code generation"}
	}
	existing, ok := c.FileServer.ReadFile(filePath)
	prompt := reasoning.Instructions + "\n\n"
	if ok {
		prompt += "Update the existing file " + filePath + ". Current content:\n" + existing + "\n\n"
	} else {
		prompt += "Create " + filePath + " with the following:\n"
	}
	prompt += step.Description + "\n\nGenerate complete, working code:"
	generated, err := c.GenerateCode(ctx, prompt)
	if err != nil {
		return StepResult{Success: false, Error: err.Error(), FilePath: filePath}
	}
	code := cleanCodeBlock(generated)
	var res FileResult
	if ok {
		res = c.FileServer.HandleModifyFile(filePath, code, false)
	} else {
		res = c.FileServer.HandleCreateFile(filePath, code)
	}
	msg := "Code generated and created: " + filePath
	if ok {
		msg = "Code generated and updated: " + filePath
	}
	return StepResult{Success: res.Success, Error: res.Error, FilePath: filePath, Message: msg}
}

// HandleTestCreation creates a test file.
func (c *Client) HandleTestCreation(ctx context.Context, step *types.Step, reasoning *types.Reasoning) StepResult {
	targetFile := step.Target
	if targetFile == "" && len(step.Dependencies) > 0 {
		targetFile = step.Dependencies[0]
	}
	var testContent string
	if reasoning != nil && reasoning.Instructions != "" {
		prompt := "Create comprehensive test cases for " + targetFile + ".\n" + reasoning.Instructions + "\n\nGenerate complete test file with multiple test cases:"
		gen, err := c.GenerateCode(ctx, prompt)
		if err == nil {
			testContent = cleanCodeBlock(gen)
		}
	}
	res := c.TestServer.HandleCreateTests(step.Target, testContent, targetFile)
	return StepResult{Success: res.Success, Error: res.Error, FilePath: res.TestPath, Message: res.Message}
}

// HandleCommandExecution runs a shell command.
func (c *Client) HandleCommandExecution(step *types.Step) StepResult {
	command := step.Target
	if command == "" {
		command = step.Command
	}
	if command == "" {
		return StepResult{Success: false, Error: "no command specified"}
	}
	cwd := c.WorkspacePath
	res := c.CommandServer.HandleExecuteCommand(command, cwd)
	return StepResult{
		Success:  res.Success,
		Error:    res.Error,
		FilePath: res.Command,
		Message:  res.Stdout,
	}
}

func cleanCodeBlock(s string) string {
	s = strings.TrimSpace(s)
	re := regexp.MustCompile("(?m)^```[\\w]*\\n?")
	s = re.ReplaceAllString(s, "")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}
