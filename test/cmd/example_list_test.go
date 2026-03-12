package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/marcuwynu23/cli-go-project-template/cmd"
)

func TestExampleListCommand(t *testing.T) {
	root := cmd.RootCmd()
	out := bytes.NewBuffer(nil)
	errBuf := bytes.NewBuffer(nil)
	root.SetOut(out)
	root.SetErr(errBuf)
	root.SetArgs([]string{"example", "list"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("example list execute: %v", err)
	}

	s := out.String()
	if s == "" {
		t.Fatal("list produced no output")
	}
	if !strings.Contains(s, "Listing") {
		t.Errorf("output should mention listing: %q", s)
	}
}
