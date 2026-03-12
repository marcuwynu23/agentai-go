package config

import (
	"os"
	"strconv"
)

// Config holds agentai configuration (file + environment).
type Config struct {
	// AI provider: gemini, openai, openrouter, ollama
	Provider      string
	APIKey        string
	Model         string
	BaseURL       string
	GeminiAPIToken string // legacy env; used when Provider is gemini and APIKey empty
	GeminiModel   string
	WorkspacePath string
	LogsPath      string
	RequestDelayMs int
	MaxRetries   int
	MCPServers   MCPServers
	LogLevel     string
	LogMaxFileSize string
	LogMaxFiles  int
}

// MCPServers holds MCP server base URLs (used for reference; we use in-process in Go).
type MCPServers struct {
	FileHandling     string
	CommandExecution string
	TestServer       string
}

// Load reads configuration: optional file (local then global) then env overrides.
// If explicitConfigPath is non-empty, only that file is used for AI config.
func Load(explicitConfigPath, cwd string) *Config {
	requestDelay, _ := strconv.Atoi(getEnv("REQUEST_DELAY", "2000"))
	maxRetries, _ := strconv.Atoi(getEnv("MAX_RETRIES", "3"))
	logMaxFiles, _ := strconv.Atoi(getEnv("LOG_MAX_FILES", "5"))
	if cwd == "" {
		cwd, _ = os.Getwd()
	}

	cfg := &Config{
		WorkspacePath:  getEnv("WORKSPACE_PATH", cwd),
		LogsPath:       getEnv("LOGS_PATH", ""),
		RequestDelayMs: requestDelay,
		MaxRetries:     maxRetries,
		MCPServers: MCPServers{
			FileHandling:     getEnv("MCP_FILE_HANDLING_URL", "http://localhost:3001"),
			CommandExecution: getEnv("MCP_COMMAND_EXECUTION_URL", "http://localhost:3002"),
			TestServer:       getEnv("MCP_TEST_SERVER_URL", "http://localhost:3003"),
		},
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		LogMaxFileSize: getEnv("LOG_MAX_FILE_SIZE", "10MB"),
		LogMaxFiles:    logMaxFiles,
	}

	// AI: file config first
	fileAI, _ := ResolveAIConfig(explicitConfigPath, cwd)
	if fileAI != nil {
		cfg.Provider = fileAI.Provider
		cfg.APIKey = fileAI.APIKey
		cfg.Model = fileAI.Model
		cfg.BaseURL = fileAI.BaseURL
	}
	// Env overrides (for backward compat and CI)
	if p := getEnv("AGENTAI_PROVIDER", ""); p != "" {
		cfg.Provider = p
	}
	if k := getEnv("GEMINI_API_TOKEN", ""); k != "" {
		if cfg.Provider == "" {
			cfg.Provider = "gemini"
		}
		if cfg.Provider == "gemini" && cfg.APIKey == "" {
			cfg.APIKey = k
		}
	}
	if k := getEnv("OPENAI_API_KEY", ""); k != "" && (cfg.Provider == "openai" || cfg.Provider == "") {
		if cfg.Provider == "" {
			cfg.Provider = "openai"
		}
		if cfg.APIKey == "" {
			cfg.APIKey = k
		}
	}
	if m := getEnv("GEMINI_MODEL", ""); m != "" {
		cfg.GeminiModel = m
	}
	if cfg.Model == "" && cfg.GeminiModel != "" {
		cfg.Model = cfg.GeminiModel
	}
	if cfg.Model == "" {
		cfg.Model = defaultModel(cfg.Provider)
	}
	if cfg.Provider == "" {
		cfg.Provider = "gemini"
	}
	return cfg
}

func defaultModel(provider string) string {
	switch provider {
	case "openai":
		return "gpt-4o-mini"
	case "openrouter":
		return "google/gemini-2.0-flash-001"
	case "ollama":
		return "llama3.2"
	default:
		return "gemini-2.5-flash"
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
