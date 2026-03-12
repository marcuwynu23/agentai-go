package cmd

import (
	"fmt"
	"os"

	"github.com/marcuwynu23/cli-go-project-template/cmd/example"
	"github.com/marcuwynu23/cli-go-project-template/cmd/version"
	"github.com/marcuwynu23/cli-go-project-template/internal/service"
	"github.com/marcuwynu23/cli-go-project-template/internal/view"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
)

// Build-time version (set via ldflags).
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "A professional CLI application template",
	Long: `A professional CLI application built with Go and Cobra.

This template demonstrates subcommands, nested commands, arguments,
global flags, and best practices for building production-ready CLIs.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Optional: validate global config, load config file, etc.
		return nil
	},
}

// Execute runs the root command and all subcommands.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// RootCmd returns the root command for testing (e.g. from test/ folder).
func RootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true // Disable the default completion command

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.app.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(version.NewCommand(func() (service.VersionProvider, *view.VersionRenderer) {
		d := deps()
		return d.VersionProvider, d.VersionView
	}))
	rootCmd.AddCommand(example.NewCommand(func() (service.ExampleUseCase, *view.ExampleRenderer, func() bool) {
		d := deps()
		return d.ExampleUseCase, d.ExampleView, func() bool { return verbose }
	}))
}
