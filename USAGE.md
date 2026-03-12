# How to Use agentai (Go)

## Quick Start

### 1. Build the CLI

```bash
go build -o agentai .
# Windows: agentai.exe
# Linux/macOS: ./agentai
```

### 2. Configure the AI provider

Config is stored in `.agentai/config.json`. Use **local** (current directory) or **global** (home directory).

**Option A: Config file (recommended)**

```bash
# Local: ./.agentai/config.json (per repository)
agentai config set provider gemini --local
agentai config set api_key YOUR_GEMINI_API_TOKEN --local
agentai config set model gemini-2.5-flash --local

# Global: ~/.agentai/config.json (all projects)
agentai config set provider gemini --global
agentai config set api_key YOUR_GEMINI_API_TOKEN --global
```

**Option B: Environment variables**

Create a `.env` file (see `.env.example`) or export:

```bash
export GEMINI_API_TOKEN=your-actual-api-token-here
# Optional
export GEMINI_MODEL=gemini-2.5-flash
```

**Show current config**

```bash
agentai config show --local    # repo config
agentai config show --global   # user config
```

### 3. Run the tool

```bash
# After building
./agentai chat "your goal here"

# Windows
agentai chat "your goal here"
```

**Examples**

```bash
agentai chat "create a hello world program"
agentai chat "create a REST API with Express.js that has a /users endpoint"
agentai chat "build a todo list application with Node.js"
```

---

## Supported providers

agentai supports four backends (raw `net/http` only; no SDK). Set **provider** and optionally **model** and **base_url**.

| Provider       | Default base URL                                              | Config keys              |
|----------------|---------------------------------------------------------------|--------------------------|
| **gemini**     | `https://generativelanguage.googleapis.com/v1beta/models`     | `api_key`, `model`, `base_url` |
| **openai**     | `https://api.openai.com/v1`                                  | `api_key`, `model`, `base_url` |
| **openrouter** | `https://openrouter.ai/api/v1`                                | `api_key`, `model`, `base_url` |
| **ollama**     | `http://localhost:11434`                                      | `model`, `base_url` (no key)   |

Defaults use the URLs above. Override with **base_url** (e.g. custom Ollama host or API proxy).

### Gemini

```bash
agentai config set provider gemini --local
agentai config set api_key YOUR_GEMINI_KEY --local
agentai config set model gemini-2.5-flash --local
# Optional: custom endpoint
agentai config set base_url https://your-proxy.com/v1beta/models --local
```

### OpenAI

```bash
agentai config set provider openai --local
agentai config set api_key sk-your-openai-key --local
agentai config set model gpt-4o-mini --local
```

### OpenRouter

```bash
agentai config set provider openrouter --local
agentai config set api_key sk-or-your-key --local
agentai config set model google/gemini-2.0-flash-001 --local
```

### Ollama (default: localhost)

```bash
agentai config set provider ollama --local
agentai config set model llama3.2 --local
# Optional: remote Ollama (default is http://localhost:11434)
agentai config set base_url http://192.168.1.55:11434 --local
```

---

## Config file locations

- **`--local`**: `<current-directory>/.agentai/config.json` (repository-scoped)
- **`--global`**: `~/.agentai/config.json` (user-scoped)

If you don’t pass `--local` or `--global`, **show** and **set** use the local path.  
When running **chat**, config is resolved in order: explicit `--config` path → local `.agentai/config.json` → global `~/.agentai/config.json` → environment variables.

---

## Understanding the output

When you run `agentai chat "goal"`:

1. **Header** – Agentic AI Code Assistant banner
2. **Project** – New project gets an AI-generated name and directory; existing project is reused
3. **Codebase analysis** – File count and any detected issues
4. **Plan** – Steps (file_creation, code_generation, test_creation, command_execution)
5. **Execution** – Per-step success/failure
6. **Completion** – “Plan execution completed!”

Example:

```
🎯 Processing goal: create a simple calculator

╭─────────────────────────────────────╮
│   🤖 Agentic AI Code Assistant       │
╰─────────────────────────────────────╯

✓ Project name: simple-calculator
╭─ Execution Plan ──────────────────────╮
│ 1. file_creation: Create package.json
│ 2. code_generation: Create calculator logic
│ 3. test_creation: Add tests
╰─────────────────────────────────────╯

[1/3] FILE_CREATION: Create package.json
  ✓ File created: package.json
...
╭─────────────────────────────────────╮
│   ✅ Plan execution completed!       │
╰─────────────────────────────────────╯
```

---

## Project structure after running

- **`<project-name>/`** – Generated project (e.g. `simple-calculator/`)
- **`<project-name>/.memory.json`** – Project state and conversation history
- **`logs/`** – Only if `LOGS_PATH` is set in the environment

---

## Commands reference

| Command | Description |
|---------|-------------|
| `agentai chat <goal>` | Generate code from a natural-language goal |
| `agentai config show [--local\|--global]` | Show current AI config |
| `agentai config set <key> <value> [--local\|--global]` | Set `provider`, `api_key`, `model`, or `base_url` |
| `agentai version` | Print version, commit, and build date |

**Global flags**

- `--config <path>` – Use a specific config file
- `-v, --verbose` – Verbose output

---

## Troubleshooting

### No API key

Set an API key via config or env:

```bash
agentai config set api_key YOUR_KEY --local
# or: export GEMINI_API_TOKEN=... or OPENAI_API_KEY=...
```

### Ollama connection refused

- Run Ollama locally (`ollama serve`) or set **base_url** to your Ollama host:
  ```bash
  agentai config set base_url http://192.168.1.55:11434 --local
  ```
- Default is `http://localhost:11434`.

### Wrong provider or model

Inspect and update config:

```bash
agentai config show --local
agentai config set provider gemini --local
agentai config set model gemini-2.5-flash --local
```

### Go version

Requires Go 1.22 or later. Check with `go version`.

---

## Example session

```bash
$ agentai config set provider ollama --local
$ agentai config set model llama3.2 --local
$ agentai chat "create a simple calculator"

🎯 Processing goal: create a simple calculator

╭─────────────────────────────────────╮
│   🤖 Agentic AI Code Assistant       │
╰─────────────────────────────────────╯

✓ Project name: simple-calculator
╭─ Execution Plan ──────────────────────╮
│ 1. file_creation: Create package.json
│ 2. code_generation: Create calculator.js
│ 3. test_creation: Add tests
│ 4. command_execution: npm install
╰─────────────────────────────────────╯

[1/4] FILE_CREATION: Create package.json
  ✓ File created: package.json
...

╭─────────────────────────────────────╮
│   ✅ Plan execution completed!       │
╰─────────────────────────────────────╯
```
