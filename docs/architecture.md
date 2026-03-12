# agentai-go Architecture

High-level layout of the agentai-go CLI and its packages.

---

## Layers

### CLI (`cmd/`)

- **root.go**: Root command `agentai`, persistent flags (`--config`, `-v`), registers `chat`, `config`, `version`.
- **chat/**: Loads config, validates API key, runs `core.ChatCommand`.
- **config/**: `config show` and `config set` for `.agentai/config.json` (--local | --global).
- **version/**: Prints version, commit, build date (from root’s build-time vars).

No business logic in `cmd/`; it delegates to `internal/`.

### Config (`internal/config`)

- **config.go**: `Load(explicitPath, cwd)` – merges file (local then global) and env (provider, api_key, model, base_url, etc.).
- **file.go**: Paths for local/global `.agentai/config.json`, `LoadAIConfig`, `SaveAIConfig`, `ResolveAIConfig`.

### Core (`internal/core`)

- **ai_core.go**: Rate limiting, retries, `GeneratePlan`, `GenerateCode`, `ReasonAboutStep`, `GenerateProjectName` (dispatches to providers).
- **gemini.go**: Gemini API HTTP client (generateContent).
- **providers.go**: Single entry point for all providers (gemini, openai, openrouter, ollama) via raw HTTP.
- **planner.go**: Builds plan prompt, calls AI, parses JSON plan into steps.
- **memory_manager.go**: Load/save `.memory.json`, conversation history, execution history.
- **codebase_analyzer.go**: Scan workspace, analyze structure, detect issues, format for AI.
- **chat_handler.go**: Orchestrates full chat: memory, project name, codebase analysis, plan, step execution via MCP client.
- **types.go**: Type aliases to `internal/types` (Plan, Step, Reasoning).

### Types (`internal/types`)

- **types.go**: `Plan`, `Step`, `Reasoning` – shared between core and mcp to avoid import cycles.

### MCP (`internal/mcp`)

- **client.go**: Coordinates file/command/test “servers”; implements step handlers (file_creation, code_generation, test_creation, command_execution). Uses injected `CodeGenFunc` from core.
- **file_server.go**: Create, modify, read files under workspace.
- **command_server.go**: Validate and execute shell commands (allow/block lists).
- **test_server.go**: Create test files (AI or template).

---

## Flow (chat)

```
User: agentai chat "build a todo app"
    ↓
cmd/chat: load config, validate API key
    ↓
core.ChatCommand: load memory, resolve project (new name or existing)
    ↓
CodebaseAnalyzer.Analyze()
    ↓
Planner.CreatePlan(goal, memory, conversation, analysis)
    ↓
For each step:
    AICore.ReasonAboutStep(step, memory)
    MCPClient.Handle*(step, reasoning)  → FileServer / CommandServer / TestServer
    MemoryManager.UpdateMemory(...)
    ↓
MemoryManager.SaveMemory()
```

---

## Tests

Tests for the removed template (example, model, service, view) have been removed. The `test/` directory is reserved for future agentai-go tests (e.g. config, providers, chat).

Run: `go test ./...`
