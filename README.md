# agentai-go – Agentic AI Code Assistant

A Go command-line tool that acts as an agentic AI code generator with planning, memory, config (local/global), and multi-provider support (Gemini, OpenAI, OpenRouter, Ollama). Same concept as the [Node.js agentai](https://github.com/marcuwynu23/agentai) implementation.

## Features

- **Goal-based code generation**: Describe what you want; agentai plans and executes it
- **Multi-provider AI**: Gemini, OpenAI, OpenRouter, or Ollama (raw HTTP; no SDK)
- **Config**: Local (`.agentai/config.json` in repo) or global (`~/.agentai/config.json`); overridable via env
- **Intelligent planning**: Breaks goals into steps (file_creation, code_generation, test_creation, command_execution)
- **Memory**: Project state and conversation in `<project-name>/.memory.json`
- **In-process “MCP”**: File, command, and test operations (no remote servers required)
- **Safe command execution**: Validation and blocklists for shell commands

## Project structure

```
.
├── main.go                 # Entry point
├── cmd/
│   ├── root.go             # Root command (agentai)
│   ├── chat/chat.go        # chat subcommand
│   ├── config/config.go    # config show / set (--local | --global)
│   └── version/version.go  # version subcommand
├── internal/
│   ├── config/             # Config load/save, .agentai paths
│   ├── core/               # AI, planner, memory, analyzer, chat handler
│   ├── types/              # Plan, Step, Reasoning
│   └── mcp/                # File, command, test servers + client
├── .env.example            # Example env (copy to .env)
├── .agentai/               # Local config (optional; create via config set)
├── go.mod
├── USAGE.md                # Usage and provider setup
├── IMPLEMENTATION.md       # Implementation summary
├── PROJECT_STRUCTURE.json  # Structure metadata
└── PROJECT_STRUCTURE.yaml
```

## Installation

1. **Go 1.22+**

2. **Clone or navigate to the project**
   ```bash
   cd agentai-go
   ```

3. **Build**
   ```bash
   go build -o agentai .
   # Windows: agentai.exe
   ```

4. **Configure** (e.g. Gemini)
   ```bash
   agentai config set provider gemini --local
   agentai config set api_key YOUR_GEMINI_API_TOKEN --local
   agentai config set model gemini-2.5-flash --local
   ```
   Or use `.env` (see `.env.example`) with `GEMINI_API_TOKEN`, etc.

## Configuration

- **Config file**: `.agentai/config.json`
  - **Local**: `<current-directory>/.agentai/config.json` (use `--local`)
  - **Global**: `~/.agentai/config.json` (use `--global`)
- **Keys**: `provider`, `api_key`, `model`, `base_url`
- **Resolution**: Explicit `--config` path → local file → global file → environment variables

Commands:

```bash
agentai config show --local     # Show repo config
agentai config show --global    # Show user config
agentai config set provider ollama --local
agentai config set model llama3.2 --local
agentai config set base_url http://192.168.1.55:11434 --local  # Optional
```

Environment (optional):

- **GEMINI_API_TOKEN** – Gemini API key (used when provider is gemini and no `api_key` in file)
- **OPENAI_API_KEY** – OpenAI key (for openai provider)
- **GEMINI_MODEL** – Model name (e.g. `gemini-2.5-flash`)
- **REQUEST_DELAY**, **MAX_RETRIES** – Rate limiting
- **WORKSPACE_PATH**, **LOGS_PATH** – Paths (defaults: cwd, empty)

### Providers and defaults

| Provider     | Default base URL                                          |
|-------------|------------------------------------------------------------|
| gemini      | `https://generativelanguage.googleapis.com/v1beta/models` |
| openai      | `https://api.openai.com/v1`                               |
| openrouter  | `https://openrouter.ai/api/v1`                            |
| ollama      | `http://localhost:11434`                                  |

Override any with `base_url` in config.

## Usage

```bash
agentai chat "your goal here"
```

Examples:

```bash
agentai chat "create a hello world program"
agentai chat "create a REST API with Express.js that has a /users endpoint"
agentai chat "build a todo list application with Node.js"
```

With Ollama (default localhost):

```bash
agentai config set provider ollama --local
agentai config set model llama3.2 --local
agentai chat "create a simple CLI in Go"
```

## Development status

- **Multi-provider AI**: Implemented (Gemini, OpenAI, OpenRouter, Ollama) via raw HTTP
- **Config**: Local/global `.agentai/config.json` and env
- **File operations**: Create, modify, read in project directory
- **Code generation**: AI-generated code with cleanup
- **Test creation**: AI-generated or template test files
- **Command execution**: Validated, safe command execution

## How it works

1. **Config**: Load provider, api_key, model, base_url (file + env).
2. **Project**: New run → AI suggests project name and directory; existing run → reuse.
3. **Analysis**: Scan codebase; summarize for planner.
4. **Plan**: AI produces a JSON plan (steps with types and dependencies).
5. **Execution**: For each step, AI reasons then the MCP client runs file/command/test logic.
6. **Memory**: Results and conversation saved to `<project-name>/.memory.json`.

## Requirements

- Go 1.22+
- API key for Gemini, OpenAI, or OpenRouter; or local Ollama (no key)

## Documentation

- **USAGE.md** – Usage, providers, config, troubleshooting
- **IMPLEMENTATION.md** – Implementation summary and architecture
- **PROJECT_STRUCTURE.json** / **PROJECT_STRUCTURE.yaml** – Structure and file roles
- **docs/** – Development guide and architecture

## License

MIT
