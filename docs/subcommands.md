# Creating Subcommands

Subcommands are top-level commands under your root CLI (e.g. `app version`, `app example`).

## Steps

1. **Create a new file** in `cmd/`, e.g. `cmd/hello.go`.

2. **Define the command** with `Use`, `Short`, and `Long`:

```go
package cmd

import "github.com/spf13/cobra"

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Say hello",
	Long:  `Print a greeting message.`,
	RunE:  runHello,
}
```

3. **Register it** in `init()`:

```go
func init() {
	rootCmd.AddCommand(helloCmd)
}
```

4. **Implement the runner**:

```go
func runHello(cmd *cobra.Command, args []string) error {
	// Your logic; return an error to exit non-zero
	return nil
}
```

## Conventions

- One file per (top-level) command keeps the tree readable.
- Use `RunE` (not `Run`) so you can return errors and get consistent exit codes.
- Keep business logic in `internal/` or testable packages; use `cmd/` for wiring and Cobra.
