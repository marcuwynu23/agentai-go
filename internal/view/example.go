package view

import (
	"fmt"
	"io"

	"github.com/marcuwynu23/cli-go-project-template/internal/model"
)

// ExampleRenderer renders example resources to the CLI.
type ExampleRenderer struct{}

// NewExampleRenderer returns a new ExampleRenderer.
func NewExampleRenderer() *ExampleRenderer {
	return &ExampleRenderer{}
}

// RenderCreated writes the "created" message to w.
func (r *ExampleRenderer) RenderCreated(w io.Writer, name string, verbose bool, force bool) {
	if verbose {
		fmt.Fprintf(w, "Creating resource %q (force=%v)\n", name, force)
	}
	fmt.Fprintf(w, "Created: %s\n", name)
}

// RenderList writes the list result to w.
func (r *ExampleRenderer) RenderList(w io.Writer, result *model.ListExampleResult, verbose bool) {
	if verbose {
		fmt.Fprintf(w, "Listing up to %d items\n", result.Total)
	}
	fmt.Fprintf(w, "Listing %d item(s)\n", result.Total)
}
