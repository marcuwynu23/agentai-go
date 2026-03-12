package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/marcuwynu23/cli-go-project-template/cmd"
)

func TestRootCommand_Help(t *testing.T) {
	root := cmd.RootCmd()
	out := bytes.NewBuffer(nil)
	errBuf := bytes.NewBuffer(nil)
	root.SetOut(out)
	root.SetErr(errBuf)
	root.SetArgs([]string{"--help"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("root --help: %v", err)
	}

	s := out.String()
	if !strings.Contains(s, "Usage:") {
		t.Errorf("help should contain Usage: %q", s)
	}
	if !strings.Contains(s, "version") {
		t.Errorf("help should list version subcommand: %q", s)
	}
	if !strings.Contains(s, "example") {
		t.Errorf("help should list example subcommand: %q", s)
	}
}
