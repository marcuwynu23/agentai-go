# CLI Development Guide

This guide explains how to extend and maintain this Go + Cobra CLI template: subcommands, nested subcommands, arguments, help, and versioning.

---

## Table of Contents

1. [Overview](#overview)
2. [Project Structure](#project-structure)
3. [Creating Subcommands](#creating-subcommands)
4. [Nested Subcommands](#nested-subcommands)
5. [Arguments](#arguments)
6. [Help and Version](#help-and-version)
7. [Testing](#testing)
8. [Building and Releasing](#building-and-releasing)

---

## Overview

This template uses [Cobra](https://github.com/spf13/cobra) for the CLI. The root command is registered in `cmd/root.go`, and each subcommand lives in its own file under `cmd/` (e.g. `cmd/version.go`, `cmd/example.go`).

- **Root command**: `app` (or your binary name)
- **Subcommands**: `version`, `example`
- **Nested under `example`**: `create`, `list`

---

## Architecture (MVC + Service Layer)

The template uses an **MVC-style layout with a service layer** to keep the codebase maintainable and avoid technical debt:

| Layer | Location | Responsibility |
|-------|----------|----------------|
| **Model** | `internal/model/` | Domain entities and DTOs (e.g. `VersionInfo`, `ExampleResource`, `CreateExampleInput`). No logic. |
| **Service** | `internal/service/` | Business logic only. Interfaces (e.g. `VersionProvider`, `ExampleUseCase`) and implementations. No I/O or formatting. |
| **View** | `internal/view/` | Presentation: render models to `io.Writer` (CLI output). No business logic. |
| **Controller** | `cmd/` | Thin: parse args/flags, call service, pass result to view. No business or output logic. |

**Why this helps:** Commands stay thin, business logic is testable without the CLI, and you can change output format or add new UIs without touching services or models.

Dependencies are wired in `cmd/deps.go` and can be replaced in tests (e.g. mock `VersionProvider`).

---

## Project Structure

```
.
├── cmd/                    # Controllers (Cobra: parse input, call service, call view)
│   ├── root.go             # Root command and persistent flags
│   ├── deps.go             # Dependency wiring (services + views; injectable for tests)
│   ├── version.go          # version subcommand
│   ├── example.go          # example parent command
│   ├── example_create.go   # example create (nested)
│   └── example_list.go     # example list (nested)
├── internal/
│   ├── model/              # Models (domain structs, input/result DTOs)
│   ├── service/           # Service layer (interfaces + implementations, business logic)
│   └── view/              # View layer (render model to io.Writer)
├── docs/
├── main.go
├── go.mod
└── README.md
```

---

## Creating Subcommands

### 1. Define the command

Create a new file under `cmd/`, e.g. `cmd/hello.go`:

```go
package cmd

import "github.com/spf13/cobra"

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Say hello",
	Long:  `Print a greeting message.`,
	RunE:  runHello,
}

func init() {
	rootCmd.AddCommand(helloCmd)
}

func runHello(cmd *cobra.Command, args []string) error {
	// Your logic here
	return nil
}
```

### 2. Register with the root

In `init()`, add:

```go
rootCmd.AddCommand(helloCmd)
```

### 3. Use and Short/Long

- **Use**: Command name and optional argument placeholders, e.g. `"hello [name]"`.
- **Short**: One-line description (shown in `app --help`).
- **Long**: Detailed description (shown in `app hello --help`).

---

## Nested Subcommands

Nested subcommands are subcommands of a non-root command (e.g. `example create`).

### 1. Create a parent command

In `cmd/example.go` we have a parent that does nothing by itself:

```go
var exampleCmd = &cobra.Command{
	Use:   "example",
	Short: "Example command with nested subcommands",
	Long:  `...`,
}

func init() {
	rootCmd.AddCommand(exampleCmd)
}
```

### 2. Add nested commands to the parent

In `cmd/example_create.go`:

```go
func init() {
	exampleCmd.AddCommand(exampleCreateCmd)  // Not rootCmd!
}
```

So the hierarchy is: `rootCmd` → `exampleCmd` → `exampleCreateCmd`.

### 3. Resulting CLI

- `app example` — shows help for the example group
- `app example create <name>` — runs the create subcommand
- `app example list` — runs the list subcommand

---

## Arguments

### Positional arguments

Use the **Args** field on the command:

| Validator           | Meaning |
|---------------------|--------|
| `cobra.NoArgs`       | No arguments allowed |
| `cobra.ExactArgs(n)` | Exactly `n` arguments |
| `cobra.MinimumNArgs(n)` | At least `n` arguments |
| `cobra.MaximumNArgs(n)` | At most `n` arguments |
| `cobra.ArbitraryArgs`  | Any number of arguments |

Example (exactly one argument):

```go
var exampleCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new example resource",
	Args:  cobra.ExactArgs(1),
	RunE:  runExampleCreate,
}

func runExampleCreate(cmd *cobra.Command, args []string) error {
	name := args[0]
	// ...
}
```

The placeholder `[name]` in **Use** is for help text only; **Args** enforces the count.

### Flags (optional “arguments”)

For optional or named inputs, use flags instead of positional args:

```go
exampleCreateCmd.Flags().BoolVarP(&createForce, "force", "f", false, "overwrite if exists")
exampleListCmd.Flags().IntVarP(&listLimit, "limit", "n", 10, "max items to show")
```

- **Persistent flags** (on `rootCmd`) apply to all subcommands (e.g. `--verbose`, `--config`).

---

## Help and Version

### Built-in help

Cobra adds automatically:

- `app --help` / `app -h` — root help
- `app version --help` — help for `version`
- `app example create --help` — help for nested `create`

Customize with **Use**, **Short**, **Long**, and **Example**:

```go
Example: "  app example create my-resource",
```

### Version command

The `version` subcommand prints version, commit, and build date. Values are set at build time with ldflags (see [Building and Releasing](#building-and-releasing)).

Version formatting lives in `internal/view`; version data comes from `internal/service`. Both can be unit tested without running the CLI.

---

## Testing

- **Model**: `internal/model` — unit tests for structs if needed.
- **Service**: `internal/service` — unit tests for business logic (no I/O).
- **View**: `internal/view` — unit tests that render to a buffer and assert output.
- **Controller (cmd)**: Execute `rootCmd` with `SetArgs()` and capture output; optionally replace `defaultDeps` with mocks (see `cmd/version_test.go`, `cmd/example_create_test.go`).

Run tests:

```bash
go test ./...
```

---

## Building and Releasing

Default version info (e.g. `dev`, `none`, `unknown`) is in `cmd/version.go`. For releases, inject real values with ldflags:

```bash
go build -ldflags "-X github.com/marcuwynu23/cli-go-project-template/cmd.Version=1.0.0 \
  -X github.com/marcuwynu23/cli-go-project-template/cmd.Commit=$(git rev-parse --short HEAD) \
  -X github.com/marcuwynu23/cli-go-project-template/cmd.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -o app .
```

Use the same module path as in your `go.mod` when setting `Version`, `Commit`, and `BuildDate`.

---

For more on Cobra, see [cobra.dev](https://cobra.dev) and [GitHub](https://github.com/spf13/cobra).
