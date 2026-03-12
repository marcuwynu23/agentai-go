package version

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// GetVersion returns version, commit, buildDate (e.g. from ldflags).
type GetVersion func() (version, commit, buildDate string)

// NewCommand returns the version subcommand.
func NewCommand(getVersion GetVersion) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Print the application version, git commit, and build date.",
		RunE: func(c *cobra.Command, args []string) error {
			v, commit, date := getVersion()
			return runVersion(c.OutOrStdout(), v, commit, date)
		},
	}
	return cmd
}

func runVersion(w io.Writer, version, commit, buildDate string) error {
	_, err := fmt.Fprintf(w, "agentai version %s\n  commit: %s\n  built:  %s\n", version, commit, buildDate)
	return err
}
