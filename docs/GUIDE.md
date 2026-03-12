# agentai-go Development Guide

This guide describes the structure of agentai-go and how to extend it.

---

## Overview

- **Binary**: `agentai`
- **Subcommands**: `chat`, `config`, `version`
- **Framework**: [Cobra](https://github.com/spf13/cobra)

---

## Project Structure

```
.
├── main.go
├── cmd/
│   ├── root.go           # Root command, flags, registers subcommands
│   ├── chat/chat.go      # chat <goal>
│   ├── config/config.go  # config show | set (--local | --global)
│   └── version/version.go
├── internal/
│   ├── config/           # .agentai/config.json load/save, paths
│   ├── core/              # AI, planner, memory, analyzer, chat handler
│   ├── types/             # Plan, Step, Reasoning (shared)
│   └── mcp/               # File, command, test servers + client
├── docs/
├── go.mod
└── README.md
```

---

## Adding a Subcommand

1. Create `cmd/<name>/<name>.go` (e.g. `cmd/hello/hello.go`).
2. Implement `func NewCommand() *cobra.Command`.
3. In `cmd/root.go` `init()`, add: `rootCmd.AddCommand(hello.NewCommand())`.

Use Cobra’s `AddCommand` and `NewCommand()` pattern as in `cmd/chat` and `cmd/config`.

---

## Key Packages

| Package           | Role |
|-------------------|------|
| `internal/config` | Config from .agentai/config.json and env |
| `internal/core`   | AI (providers), planner, memory, codebase analyzer, chat flow |
| `internal/types`  | Shared structs (Plan, Step, Reasoning) |
| `internal/mcp`    | File/command/test operations and client |

---

## Building and Releasing

```bash
make build          # Binary in dist/agentai (or .exe)
make release-build  # With version/commit/date ldflags
make test           # go test ./...
```

Version is set via ldflags; see `Makefile` and `cmd/root.go` (Version, Commit, BuildDate).

---

## Documentation

- **README.md** – Overview and quick start
- **USAGE.md** – Config, providers, commands, troubleshooting
- **IMPLEMENTATION.md** – Implementation summary and architecture
- **PROJECT_STRUCTURE.json / .yaml** – File and directory roles
