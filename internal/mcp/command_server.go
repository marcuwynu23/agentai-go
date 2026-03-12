package mcp

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// CommandResult is the result of command execution.
type CommandResult struct {
	Success  bool   `json:"success"`
	Command  string `json:"command,omitempty"`
	Stdout   string `json:"stdout,omitempty"`
	Stderr   string `json:"stderr,omitempty"`
	Error    string `json:"error,omitempty"`
	ExitCode int    `json:"exitCode,omitempty"`
}

// CommandServer runs allowed shell commands with validation.
type CommandServer struct {
	WorkspacePath   string
	AllowedCommands []string
	BlockedCommands []string
}

// NewCommandServer creates a CommandServer.
func NewCommandServer(workspacePath string) *CommandServer {
	if workspacePath == "" {
		workspacePath, _ = os.Getwd()
	}
	return &CommandServer{
		WorkspacePath:   workspacePath,
		AllowedCommands: []string{"npm", "node", "git", "go", "go mod"},
		BlockedCommands: []string{"rm", "del", "format", "fdisk"},
	}
}

// ValidateCommand checks if the command is allowed and not dangerous.
func (c *CommandServer) ValidateCommand(command string) (valid bool, reason string) {
	parts := strings.Fields(strings.TrimSpace(command))
	if len(parts) == 0 {
		return false, "empty command"
	}
	base := strings.ToLower(parts[0])
	for _, b := range c.BlockedCommands {
		if strings.Contains(base, b) {
			return false, "command '" + base + "' is blocked for safety"
		}
	}
	allowed := false
	for _, a := range c.AllowedCommands {
		if strings.HasPrefix(base, a) || strings.Contains(command, a) {
			allowed = true
			break
		}
	}
	if !allowed {
		return false, "command '" + base + "' is not in allowed list"
	}
	dangerous := []*regexp.Regexp{
		regexp.MustCompile(`rm\s+-rf`),
		regexp.MustCompile(`del\s+/`),
		regexp.MustCompile(`>\s*/dev`),
		regexp.MustCompile(`&&\s*rm`),
		regexp.MustCompile(`\|\s*sh\s*$`),
		regexp.MustCompile(`eval\s*\(`),
	}
	for _, re := range dangerous {
		if re.MatchString(command) {
			return false, "command contains dangerous patterns"
		}
	}
	return true, ""
}

// HandleExecuteCommand runs the command in WorkspacePath.
func (c *CommandServer) HandleExecuteCommand(command, cwd string) CommandResult {
	valid, reason := c.ValidateCommand(command)
	if !valid {
		return CommandResult{Success: false, Command: command, Error: reason}
	}
	workDir := cwd
	if workDir == "" {
		workDir = c.WorkspacePath
	}
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = workDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		exitCode := 1
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		}
		return CommandResult{
			Success:  false,
			Command:  command,
			Stdout:   string(out),
			Stderr:   string(out),
			Error:    err.Error(),
			ExitCode: exitCode,
		}
	}
	return CommandResult{
		Success: true,
		Command: command,
		Stdout:  strings.TrimSpace(string(out)),
		Stderr:  "",
	}
}
