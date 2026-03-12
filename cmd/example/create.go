package example

import (
	"github.com/marcuwynu23/cli-go-project-template/internal/model"
	"github.com/spf13/cobra"
)

var createForce bool

func newCreateCmd(getDeps GetDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new example resource",
		Long: `Create a new example resource by name.

The name argument is required. Use --force to overwrite existing resources.`,
		Args:    cobra.ExactArgs(1),
		Example: "  app example create my-resource",
		RunE:    runCreate(getDeps),
	}
	cmd.Flags().BoolVarP(&createForce, "force", "f", false, "overwrite if resource exists")
	return cmd
}

func runCreate(getDeps GetDeps) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, args []string) error {
		uc, v, getVerbose := getDeps()
		name := args[0]
		res, err := uc.Create(model.CreateExampleInput{Name: name, Force: createForce})
		if err != nil {
			return err
		}
		v.RenderCreated(c.OutOrStdout(), res.Name, getVerbose(), createForce)
		return nil
	}
}
