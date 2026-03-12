package view_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/marcuwynu23/cli-go-project-template/internal/model"
	"github.com/marcuwynu23/cli-go-project-template/internal/view"
)

func TestExampleRenderer_RenderCreated(t *testing.T) {
	r := view.NewExampleRenderer()
	var buf bytes.Buffer
	r.RenderCreated(&buf, "my-resource", false, false)
	s := buf.String()
	if !strings.Contains(s, "Created: my-resource") {
		t.Errorf("output missing Created: %q", s)
	}
}

func TestExampleRenderer_RenderCreated_Verbose(t *testing.T) {
	r := view.NewExampleRenderer()
	var buf bytes.Buffer
	r.RenderCreated(&buf, "x", true, true)
	s := buf.String()
	if !strings.Contains(s, "Creating resource \"x\" (force=true)") {
		t.Errorf("verbose output missing: %q", s)
	}
}

func TestExampleRenderer_RenderList(t *testing.T) {
	r := view.NewExampleRenderer()
	var buf bytes.Buffer
	r.RenderList(&buf, &model.ListExampleResult{Total: 10}, false)
	s := buf.String()
	if !strings.Contains(s, "Listing 10 item(s)") {
		t.Errorf("output missing count: %q", s)
	}
}
