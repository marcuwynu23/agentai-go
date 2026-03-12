package example

import (
	"github.com/marcuwynu23/cli-go-project-template/internal/service"
	"github.com/marcuwynu23/cli-go-project-template/internal/view"
	"github.com/spf13/cobra"
)

// GetDeps returns ExampleUseCase, ExampleRenderer, and getVerbose (called at runtime so tests can inject).
type GetDeps func() (service.ExampleUseCase, *view.ExampleRenderer, func() bool)

// NewCommand returns the example parent command with create and list subcommands. getDeps is called at run time.
func NewCommand(getDeps GetDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "example",
		Short: "Example command with nested subcommands",
		Long: `Manage example resources. Use the nested subcommands to create or list items.

Examples:
  app example create my-item
  app example list
  app example list --limit 10`,
	}
	cmd.AddCommand(newCreateCmd(getDeps))
	cmd.AddCommand(newListCmd(getDeps))
	return cmd
}
