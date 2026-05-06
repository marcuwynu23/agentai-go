# Contributing to AgentAI

Thanks for your interest in contributing to AgentAI! We welcome contributions of all kinds: bug fixes, features, documentation, and suggestions.

AgentAI is an AI-powered code assistant that helps developers write code faster with intelligent planning and multi-provider AI support.

---

## Getting Started

1. **Fork the repository**
   ```bash
   # Fork on GitHub, then clone your fork
   git clone https://github.com/YOUR_USERNAME/agentai.git
   cd agentai
   ```

2. **Add upstream remote**
   ```bash
   git remote add upstream https://github.com/marcuwynu23/agentai.git
   ```

3. **Create a new branch from `main`**
   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Install dependencies and set up project**
   ```bash
   # Install Go dependencies
   go mod tidy
   go mod download
   
   # Build the project
   go build -o agentai .
   
   # Test the build
   ./agentai --help
   ```

5. **Set up your AI provider configuration**
   ```bash
   # Copy example environment file
   cp .env.example .env
   
   # Edit .env with your API keys
   # Or use: agentai config set provider gemini --local
   # agentai config set api_key YOUR_API_KEY --local
   ```

---

## Branching Strategy

We follow a structured branching approach:

### Main Branches

- `main` → Production-ready code
- `develop` → Integration branch for ongoing development

### Supporting Branches

Use the following naming conventions:

- `feature/<short-description>` → New features
- `fix/<short-description>` → Bug fixes
- `chore/<short-description>` → Maintenance tasks
- `docs/<short-description>` → Documentation updates
- `refactor/<short-description>` → Code improvements without behavior change
- `test/<short-description>` → Adding or updating tests

### AgentAI-Specific Examples

```
feature/add-cloudflare-provider
fix/ollama-connection-error
docs/update-ai-provider-setup
refactor/improve-planning-algorithm
test/add-provider-unit-tests
chore/update-dependencies
```

---

## Development Workflow

1. **Create a branch from `main`** (AgentAI uses main as primary development branch)
2. **Make your changes in a focused branch**
3. **Follow AgentAI's coding style and conventions**
4. **Add or update tests when applicable**
5. **Run local checks before submitting**:

```bash
# Run tests
go test ./...

# Build the project
go build -o agentai .

# Test with different providers
./agentai config set provider gemini --local
./agentai chat --help

# Format code
go fmt ./...

# Lint (optional)
golangci-lint run
```

---

## Commit Messages (Conventional Commits)

We follow the **Conventional Commits** specification.

### Format

```
<type>(optional scope): <short description>
```

### Common Types

- `feat` → New feature
- `fix` → Bug fix
- `docs` → Documentation changes
- `style` → Formatting (no code logic changes)
- `refactor` → Code restructuring
- `test` → Adding/updating tests
- `chore` → Maintenance

### AgentAI Examples

```
feat(ai): add claude provider support
fix(providers): handle timeout errors in openai calls
docs(readme): update cloudflare setup instructions
refactor(planner): improve step dependency resolution
test(core): add integration tests for memory manager
chore(deps): update bubbletea to latest version
```

### Rules

- Use lowercase for type and description
- Keep messages concise and meaningful
- Use the body for additional context if needed

---

## Pull Request Process

1. **Ensure your branch is up to date with `main`**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Verify all tests and checks pass**
   ```bash
   go test ./...
   go build -o agentai .
   ```

3. **Open a pull request targeting `main`**

4. **Clearly describe your changes**:
   - What changed (specific features/fixes)
   - Why it was needed (problem statement)
   - How to test your changes
   - Any breaking changes

5. **Test with multiple AI providers** (if applicable):
   - Test with at least one free provider (Ollama if available)
   - Verify configuration changes work correctly

6. **Use PR template** if available:
   - [`.github/PULL_REQUEST_TEMPLATE.md`](.github/PULL_REQUEST_TEMPLATE.md)

**Optional but recommended**:
- Include screenshots of the TUI interface if UI changes
- Add logs or examples of AI interactions
- Document any new configuration options

---

## Reporting Issues

When reporting bugs, please use the provided template:

- [`.github/ISSUE_TEMPLATE/bug_report.md`](.github/ISSUE_TEMPLATE/bug_report.md)

**For AgentAI-specific issues, please include**:

- **AI Provider**: Which AI provider you're using (Gemini, OpenAI, Ollama, etc.)
- **Configuration**: Your provider configuration (without API keys)
- **Goal/Command**: The specific goal or command that caused the issue
- **Error Logs**: Full error output from the terminal
- **Environment**: OS, Go version, and any relevant system details
- **Expected vs Actual**: What you expected to happen vs what actually happened

---

## Suggestions & Feature Requests

For feature requests and suggestions, please use:

- [`.github/ISSUE_TEMPLATE/feature_request.md`](.github/ISSUE_TEMPLATE/feature_request.md)

**For AgentAI feature requests, please include**:

- **Use Case**: Describe the specific development scenario
- **AI Provider**: Which provider this feature would apply to
- **Proposed Solution**: How you envision the feature working
- **Alternatives**: Any other approaches you've considered
- **Examples**: Sample commands or workflows that would benefit

**Popular feature request areas**:
- New AI provider support
- Enhanced planning algorithms
- Better memory management
- Improved TUI/UX
- Additional file operations

---

## Code of Conduct

This project follows the guidelines defined in:

- [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md)

Be respectful and constructive in all interactions.
Harassment or inappropriate behavior will not be tolerated.

---

## AgentAI-Specific Guidelines

### Code Style
- Follow Go conventions and idioms
- Use meaningful variable and function names
- Add comments for complex AI logic
- Keep provider implementations consistent

### Testing
- Test provider implementations with mock responses
- Test configuration parsing and validation
- Test error handling for network issues
- Include integration tests when possible

### AI Provider Development
- Follow the existing provider pattern in `internal/core/providers.go`
- Implement proper error handling and retries
- Support the standard configuration fields
- Add documentation for new providers

### Memory and Planning
- Maintain backward compatibility for memory format
- Test planning edge cases and error scenarios
- Consider performance for large codebases

---

## Notes

- Maintainers may request changes before merging
- Not all contributions may be accepted, but all will be reviewed
- Focus on user experience and reliability
- Consider the impact on existing users when making changes

---

Thanks again for contributing to AgentAI! Your contributions help make AI-powered development more accessible to everyone.
