# Architecture: MVC + Service Layer

This template uses an **MVC-style layout with a dedicated service layer** so the CLI stays maintainable and does not accumulate technical debt as it grows.

---

## Layers

### Model (`internal/model/`)

- **Purpose**: Domain data only — structs, input DTOs, result DTOs.
- **No**: I/O, formatting, or business rules.
- **Examples**: `VersionInfo`, `ExampleResource`, `CreateExampleInput`, `ListExampleResult`.

### Service (`internal/service/`)

- **Purpose**: All business logic. Interfaces define use cases; implementations contain the rules.
- **No**: Direct I/O (no `fmt` to stdout), no knowledge of CLI or output format.
- **Examples**: `VersionProvider.GetVersion()`, `ExampleUseCase.Create()`, `ExampleUseCase.List()`.
- **Testing**: Unit test with in-memory or mock dependencies; no CLI needed.

### View (`internal/view/`)

- **Purpose**: Presentation only. Take model data and write to an `io.Writer` (e.g. stdout).
- **No**: Business logic or decisions.
- **Examples**: `VersionRenderer.Render(w, info)`, `ExampleRenderer.RenderCreated()`, `RenderList()`.
- **Testing**: Call render with a buffer and assert the string output.

### Controller (`cmd/`)

- **Purpose**: Thin glue. Parse args/flags, call service, pass results to view, handle errors.
- **No**: Business logic or output formatting.
- **Flow**: `args/flags` → **service** → **model** → **view** → `io.Writer`.
- **Structure**: One folder per command group. Root command and wiring live in `cmd/`; each subcommand (or group) has its own package under `cmd/`:
  - `cmd/root.go` — root command, build version vars, persistent flags; wires subcommands in `init()`.
  - `cmd/deps.go` — `Deps` struct and `deps()`; dependencies are created lazily and can be reset in tests via `ResetDepsForTest()`.
  - `cmd/version/` — version subcommand; `NewCommand(getDeps)` so deps are read at **run time** (test-friendly).
  - `cmd/example/` — example parent and nested `create` / `list`; same run-time `getDeps` pattern.
- **Test hooks**: `RootCmd()` returns the root command for tests; `ResetDepsForTest()` clears cached deps so the next run uses current `Version`/`Commit`/`BuildDate`.

---

## Test layout (`test/`)

All tests live under a **root-level `test/`** directory, mirroring the source layout:

- **`test/cmd/`** — CLI (controller) tests. Package `cmd_test`; uses `cmd.RootCmd()`, `cmd.ResetDepsForTest()`, and runs the full command tree.
- **`test/internal/model/`** — Model tests. Package `model_test`; tests `internal/model` via exported types only.
- **`test/internal/service/`** — Service tests. Package `service_test`; tests `internal/service` (e.g. `NewVersionService`, `NewExampleService`, `ErrInvalidInput`).
- **`test/internal/view/`** — View tests. Package `view_test`; tests renderers by writing to a buffer and asserting output.

Run everything with `go test ./...` from the repo root.

---

## Dependency flow

```
User input (args, flags)
        ↓
   cmd (controller)
        ↓
  service (business logic)  →  model (data)
        ↓
   view (presentation)
        ↓
   io.Writer (stdout)
```

---

## Why this reduces tech debt

1. **Single responsibility** — Each layer has one job; changes are localized.
2. **Testability** — Services and views are tested without running the CLI.
3. **Swappable UI** — New output formats or UIs only touch view (and maybe cmd); services stay the same.
4. **Clear boundaries** — No business logic in cmd, no I/O in service, no logic in view.
5. **Injectible deps** — `cmd/deps.go` and run-time `getDeps` in subcommands allow fresh deps each run and easy test overrides (e.g. `ResetDepsForTest()`).
6. **Centralized tests** — All tests under `test/` keep source trees clean and make it obvious where to add or run tests.

When adding a new feature: add or extend **models** → implement or extend **service** → add **view** renderers → add a **command** under `cmd/` (or new `cmd/<name>/`) and wire it in `cmd/root.go`.
