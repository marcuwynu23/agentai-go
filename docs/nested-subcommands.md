# Nested Subcommands

Nested subcommands are commands under another command (e.g. `app example create`, `app example list`).

## Parent command

Create a **parent** command that groups nested commands. It usually has no `Run`/`RunE` and only exists to hold children:

```go
// cmd/example.go
var exampleCmd = &cobra.Command{
	Use:   "example",
	Short: "Example command with nested subcommands",
	Long:  `...`,
}

func init() {
	rootCmd.AddCommand(exampleCmd)
}
```

## Child (nested) command

Create the nested command and **add it to the parent**, not to `rootCmd`:

```go
// cmd/example_create.go
var exampleCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new resource",
	RunE:  runExampleCreate,
}

func init() {
	exampleCmd.AddCommand(exampleCreateCmd)  // Parent is exampleCmd
}
```

## Resulting hierarchy

- `app` — root
- `app example` — parent (shows help for example)
- `app example create` — nested
- `app example list` — nested

Use one file per nested command (e.g. `example_create.go`, `example_list.go`) for clarity.
