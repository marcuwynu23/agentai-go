# AgentAI v1.0.0 Release Notes

**Release Date:** May 6, 2026

## Overview

This is the initial release of AgentAI, an AI-powered code assistant with a modern Terminal User Interface (TUI) that helps developers write code faster with intelligent planning and multi-provider AI support.

## What's New

### Core Features
- **AI-Powered Code Generation**: Generate complete projects from natural language descriptions
- **Multi-Provider Support**: Support for Gemini, OpenAI, OpenRouter, Ollama, and Cloudflare AI Gateway
- **Interactive TUI**: Modern terminal interface with real-time progress updates
- **Smart Planning**: AI-driven step-by-step project creation
- **Memory Management**: Persistent conversation and project state
- **Codebase Analysis**: Automatic project structure analysis and issue detection

### Supported Platforms
- **Linux**: amd64, arm64
- **Windows**: amd64, arm64  
- **macOS**: amd64, arm64

### AI Providers
- **Gemini**: Google's generative AI models
- **OpenAI**: GPT models with strong capabilities
- **OpenRouter**: Access to multiple AI models through one service
- **Ollama**: Local model hosting
- **Cloudflare AI Gateway**: Enterprise AI with analytics

## Getting Started

### Installation
```bash
# Clone the repository
git clone https://github.com/marcuwynu23/agentai.git
cd agentai

# Build the application
go build -o agentai .

# Configure your AI provider
agentai config set provider gemini --local
agentai config set api_key YOUR_API_KEY --local
agentai config set model gemini-2.5-flash --local

# Launch the TUI
./agentai chat
```

### First Use
1. Launch AgentAI with `./agentai chat`
2. Type your development goal in natural language
3. Watch as AgentAI plans and executes your project
4. Navigate through conversation history
5. Exit cleanly with Ctrl+C

## Configuration

### Local vs Global Settings
- **Local**: Project-specific configuration in `.agentai/config.json`
- **Global**: User-wide configuration in `~/.agentai/config.json`
- **Environment**: Override settings with environment variables

### Supported Environment Variables
- `GEMINI_API_TOKEN`: Gemini API key
- `OPENAI_API_KEY`: OpenAI API key
- `REQUEST_DELAY`: Rate limiting delay
- `MAX_RETRIES`: Maximum retry attempts
- `WORKSPACE_PATH`: Custom workspace directory
- `LOGS_PATH`: Custom logs directory

## Usage Examples

### Basic Projects
```
create a hello world program
build a REST API with Express.js
create a todo list application with React
```

### Advanced Features
- **Project Analysis**: Automatic codebase understanding
- **Step-by-Step Execution**: File creation, code generation, testing, command execution
- **Real-time Progress**: Live updates during project creation
- **Memory Persistence**: Resume conversations and maintain context

## Technical Details

### Architecture
- **Language**: Go 1.22+
- **TUI Framework**: Bubbletea for modern terminal interfaces
- **AI Integration**: Direct HTTP clients (no SDK dependencies)
- **Configuration**: JSON-based with environment variable support
- **Memory**: JSON-based conversation and state management

### Project Structure
- Generated projects include:
  - Source code files
  - Package configuration (package.json, go.mod, etc.)
  - Test files
  - Project memory (.memory.json)

## Known Limitations

### Current Limitations
- Requires internet connection for AI providers
- Local Ollama setup required for offline usage
- Large projects may require multiple iterations
- Some complex architectures may need manual refinement

### Platform Notes
- **Windows**: Use `agentai.exe` after building
- **Linux/macOS**: Use `./agentai` after building
- **Docker**: Not yet supported (planned for future releases)

## Documentation

- **[USAGE.md](../../USAGE.md)**: Complete usage guide
- **[CONTRIBUTING.md](../../CONTRIBUTING.md)**: Development guidelines
- **[ARCHITECTURE.md](../architecture.md)**: Technical architecture
- **[README.md](../../README.md)**: Project overview and quick start

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](../../CONTRIBUTING.md) for:
- Code style and conventions
- Pull request process
- Issue reporting
- Feature requests

## License

This release is licensed under the MIT License. See [LICENSE](../../LICENSE) for full details.

---

## What's Next

### Planned for v1.1.0
- Enhanced error handling and recovery
- Additional AI provider support
- Docker containerization
- Plugin system for custom providers
- Improved project templates

### Long-term Roadmap
- Web interface
- Team collaboration features
- Advanced project analysis
- Integration with popular IDEs

---

**Thank you for using AgentAI!**

For questions, issues, or suggestions, please visit our [GitHub repository](https://github.com/marcuwynu23/agentai).
