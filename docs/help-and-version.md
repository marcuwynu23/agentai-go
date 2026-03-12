# Help and Version

## Help

Cobra adds help automatically:

- `app --help` or `app -h` — root command help
- `app version --help` — help for `version`
- `app example create --help` — help for nested `create`

### Customizing help

Set these on any `*cobra.Command`:

- **Use**: Command name and argument placeholders, e.g. `"create [name]"`.
- **Short**: One-line description (shown in parent’s help list).
- **Long**: Full description (shown in this command’s `--help`).
- **Example**: Example usage line(s):

```go
Example: "  app example create my-resource",
```

## Version

The template includes a **version** subcommand that prints:

- Version (e.g. `1.0.0` or `dev`)
- Git commit (e.g. `abc1234`)
- Build date (e.g. `2024-01-15T12:00:00Z`)

Default values are in `cmd/version.go`. For releases, override them at build time with ldflags (see README or [GUIDE.md](GUIDE.md#building-and-releasing)).

Formatting lives in `internal/version` so it can be unit tested without running the CLI.
