package view

import (
	"fmt"
	"io"

	"github.com/marcuwynu23/cli-go-project-template/internal/model"
)

// VersionRenderer renders version info to the CLI.
type VersionRenderer struct{}

// NewVersionRenderer returns a new VersionRenderer.
func NewVersionRenderer() *VersionRenderer {
	return &VersionRenderer{}
}

// Render writes the version info to w (view layer: presentation only).
func (r *VersionRenderer) Render(w io.Writer, info model.VersionInfo) {
	fmt.Fprintf(w, "version %s\n", info.Version)
	fmt.Fprintf(w, "  commit: %s\n", info.Commit)
	fmt.Fprintf(w, "  built:  %s\n", info.BuildDate)
}
