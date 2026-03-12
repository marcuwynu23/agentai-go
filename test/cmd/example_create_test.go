package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/marcuwynu23/cli-go-project-template/cmd"
)

func TestExampleCreateCommand(t *testing.T) {
	root := cmd.RootCmd()
	out := bytes.NewBuffer(nil)
	errBuf := bytes.NewBuffer(nil)
	root.SetOut(out)
	root.SetErr(errBuf)
	root.SetArgs([]string{"example", "create", "my-resource"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("example create execute: %v", err)
	}

	s := out.String()
	if !strings.Contains(s, "my-resource") {
		t.Errorf("output should contain resource name: %q", s)
	}
}

func TestExampleCreateCommand_RequiresArg(t *testing.T) {
	root := cmd.RootCmd()
	errBuf := bytes.NewBuffer(nil)
	root.SetOut(nil)
	root.SetErr(errBuf)
	root.SetArgs([]string{"example", "create"})

	err := root.Execute()
	if err == nil {
		t.Fatal("expected error when no args provided")
	}
}
