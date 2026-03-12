package types

// Plan is the result of planning: goal + steps.
type Plan struct {
	Goal  string  `json:"goal"`
	Steps []*Step `json:"steps"`
}

// Step is a single step in the execution plan.
type Step struct {
	ID           string   `json:"id"`
	Type         string   `json:"type"`
	Description  string   `json:"description"`
	Target       string   `json:"target"`
	Dependencies []string `json:"dependencies"`
	Content      string   `json:"content,omitempty"`
	Command      string   `json:"command,omitempty"`
}

// Reasoning holds AI reasoning for a step.
type Reasoning struct {
	Approach       string   `json:"approach"`
	Instructions   string   `json:"instructions"`
	Considerations []string `json:"considerations"`
}
