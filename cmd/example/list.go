package example

import (
	"github.com/marcuwynu23/cli-go-project-template/internal/model"
	"github.com/spf13/cobra"
)

var listLimit int
var listAll bool

func newListCmd(getDeps GetDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List example resources",
		Long:  `List all example resources. Use --limit to cap the number of items, or --all to show everything.`,
		Args:  cobra.NoArgs,
		RunE:  runList(getDeps),
	}
	cmd.Flags().IntVarP(&listLimit, "limit", "n", 10, "maximum number of items to show")
	cmd.Flags().BoolVarP(&listAll, "all", "a", false, "show all items (ignores limit)")
	return cmd
}

func runList(getDeps GetDeps) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, args []string) error {
		uc, v, getVerbose := getDeps()
		result, err := uc.List(model.ListExampleInput{Limit: listLimit, All: listAll})
		if err != nil {
			return err
		}
		v.RenderList(c.OutOrStdout(), result, getVerbose())
		return nil
	}
}
