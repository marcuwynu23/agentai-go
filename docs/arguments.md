# Arguments and Flags

## Positional arguments

Control how many positional arguments a command accepts with the **Args** field:

| Validator             | Effect |
|-----------------------|--------|
| `cobra.NoArgs`        | No arguments allowed |
| `cobra.ExactArgs(1)`  | Exactly one argument |
| `cobra.MinimumNArgs(1)` | At least one |
| `cobra.MaximumNArgs(1)` | At most one |
| `cobra.ArbitraryArgs` | Any number |

Example:

```go
var exampleCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Args:  cobra.ExactArgs(1),
	RunE:  runExampleCreate,
}

func runExampleCreate(cmd *cobra.Command, args []string) error {
	name := args[0]
	// ...
}
```

The `[name]` in **Use** is for help only; **Args** enforces the count. Cobra will print an error and usage if the user passes the wrong number of arguments.

## Flags

For optional or named inputs, use flags.

**Local flags** (only on this command):

```go
exampleCreateCmd.Flags().BoolVarP(&createForce, "force", "f", false, "overwrite if exists")
exampleListCmd.Flags().IntVarP(&listLimit, "limit", "n", 10, "max items")
```

**Persistent flags** (root; apply to all commands):

```go
rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
```

- `VarP` = variant with short flag (e.g. `-f`, `-n`).
- Use `Flags()` for command-specific, `PersistentFlags()` on root for global behavior.
