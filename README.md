# CLI Go Project Template

A **console/cli project template** for building command-line applications in Go with [Cobra](https://github.com/spf13/cobra). Use it as a starting point for tools, automation, or any CLI product.

---

## Features

| Feature                 | Description                                                            |
| ----------------------- | ---------------------------------------------------------------------- |
| **Cobra**               | Industry-standard CLI framework (used by Kubernetes, Hugo, GitHub CLI) |
| **Subcommands**         | Top-level commands (e.g. `version`, `example`)                         |
| **Nested subcommands**  | Command groups (e.g. `example create`, `example list`)                 |
| **Arguments & flags**   | Positional args with validation, local and persistent flags            |
| **Help & version**      | Auto-generated help and a version command with build-time info         |
| **MVC + service layer** | Model / service / view / controller split to avoid tech debt           |
| **Tests**               | Unit tests for model, service, view; command tests for CLI             |
| **Docs**                | In-repo guide for extending the template                               |

---

## Quick start

### Prerequisites

- **Go 1.21+**

### Build and run

```bash
# Clone or use this repo as template
cd cli-go-project-template

# Download dependencies
go mod download

# Build
go build -o app .

# Run
./app --help
./app version
./app example create my-item
./app example list --limit 5
```

### Using Make (optional)

```bash
make build      # Build binary (version=dev)
make test       # Run tests
make version-build   # Build with git version/commit/date
make clean      # Remove build artifacts
```

---

## Architecture (MVC + service layer)

The project is structured to keep **controllers thin** and **business logic in services**, so it stays maintainable as it grows:

| Layer          | Directory           | Role                                                                  |
| -------------- | ------------------- | --------------------------------------------------------------------- |
| **Model**      | `internal/model/`   | Domain structs and input/result DTOs (no logic).                      |
| **Service**    | `internal/service/` | Business logic; interfaces + implementations (no I/O, no formatting). |
| **View**       | `internal/view/`    | Render models to CLI output (`io.Writer`).                            |
| **Controller** | `cmd/`              | Parse args/flags, call service, call view.                            |

Dependencies are wired in `cmd/deps.go` and can be swapped in tests (e.g. mock `VersionProvider` or `ExampleUseCase`).

---

## Project structure

```
.
├── cmd/                      # Controllers (Cobra: input → service → view)
│   ├── root.go               # Root command, persistent flags
│   ├── deps.go               # Service/view wiring (injectable for tests)
│   ├── version.go            # version subcommand
│   ├── example.go            # example parent command
│   ├── example_create.go     # example create (nested)
│   ├── example_list.go       # example list (nested)
│   └── *_test.go             # Command tests
├── internal/
│   ├── model/                # Domain entities, input/result DTOs
│   ├── service/              # Business logic (interfaces + impl)
│   └── view/                 # CLI presentation (render to io.Writer)
├── docs/
│   ├── GUIDE.md              # Full guide (incl. architecture)
│   ├── subcommands.md        # Adding subcommands
│   ├── nested-subcommands.md # Nested command groups
│   ├── arguments.md          # Args and flags
│   └── help-and-version.md   # Help and version command
├── main.go
├── go.mod
├── Makefile
└── README.md
```

---

## Command overview

| Command                     | Description                                               |
| --------------------------- | --------------------------------------------------------- |
| `app`                       | Root; show help                                           |
| `app --help`, `app -h`      | Global help                                               |
| `app version`               | Print version, commit, build date                         |
| `app example`               | Example parent; show nested commands                      |
| `app example create <name>` | Create a resource (e.g. `app example create my-resource`) |
| `app example list`          | List resources (optional: `--limit`, `--all`)             |

**Global flags** (available on all commands):

- `--config` — config file path
- `-v`, `--verbose` — verbose output

---

## Version and release builds

Default build uses placeholder version info (`dev`, `none`, `unknown`). For releases, inject real values at build time:

```bash
go build -ldflags "\
  -X github.com/marcuwynu23/cli-go-project-template/cmd.Version=1.0.0 \
  -X github.com/marcuwynu23/cli-go-project-template/cmd.Commit=$(git rev-parse --short HEAD) \
  -X github.com/marcuwynu23/cli-go-project-template/cmd.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -o app .
```

Or use the Makefile:

```bash
VERSION=1.0.0 make version-build
```

---

## Testing

Tests follow the layers: unit tests for **model**, **service**, and **view**; command tests for the **CLI** (controllers).

```bash
go test ./...
```

- **`internal/model`** — model structs
- **`internal/service`** — business logic (e.g. `Create`, `List`, `GetVersion`)
- **`internal/view`** — output formatting (render to buffer, assert content)
- **`cmd/*_test.go`** — run root with `SetArgs`, capture output; optionally replace `defaultDeps` with mocks

---

## Documentation

Detailed guides are in the **`docs/`** folder:

| Document                                                | Contents                                                                                                   |
| ------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| [**GUIDE.md**](docs/GUIDE.md)                           | Full guide: architecture, structure, subcommands, nested commands, args, help, version, testing, releasing |
| [**architecture.md**](docs/architecture.md)             | MVC + service layer: model, service, view, controller and dependency flow                                  |
| [**subcommands.md**](docs/subcommands.md)               | How to add a new top-level subcommand                                                                      |
| [**nested-subcommands.md**](docs/nested-subcommands.md) | How to add nested commands under a parent                                                                  |
| [**arguments.md**](docs/arguments.md)                   | Positional arguments and flags                                                                             |
| [**help-and-version.md**](docs/help-and-version.md)     | Customizing help and the version command                                                                   |

Start with **GUIDE.md** when extending or customizing the template.

---

## Customizing the template

1. **Rename the module**  
   Update `go.mod` and replace `github.com/marcuwynu23/cli-go-project-template` in imports and ldflags.

2. **Rename the binary**  
   Change the root command `Use` in `cmd/root.go` (e.g. from `app` to `mycli`).

3. **Add subcommands**  
   Follow [docs/subcommands.md](docs/subcommands.md) and [docs/nested-subcommands.md](docs/nested-subcommands.md).

4. **Add arguments and flags**  
   Use [docs/arguments.md](docs/arguments.md) and Cobra’s `Args` and `Flags()`.

5. **Replace example commands**  
   Remove or replace `cmd/example*.go` with your own commands.

---

## License

Use this template freely for personal or commercial projects. Consider a mention or star if it helps you ship.

---

**Built with [Cobra](https://github.com/spf13/cobra) · Go CLI template**
