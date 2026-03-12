package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// AIConfig is the AI provider config stored in .agentai/config.json.
type AIConfig struct {
	Provider string `json:"provider"` // gemini, openai, openrouter, ollama
	APIKey   string `json:"api_key"`
	Model    string `json:"model"`
	BaseURL  string `json:"base_url,omitempty"`
}

// ConfigDir is the name of the config directory under repo or home.
const ConfigDir = ".agentai"

// ConfigFileName is the config file name.
const ConfigFileName = "config.json"

// LocalPath returns path to local config: <cwd>/.agentai/config.json.
func LocalPath(cwd string) string {
	if cwd == "" {
		cwd, _ = os.Getwd()
	}
	return filepath.Join(cwd, ConfigDir, ConfigFileName)
}

// GlobalPath returns path to global config: $HOME/.agentai/config.json.
func GlobalPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ConfigDir, ConfigFileName)
}

// LoadAIConfig reads AIConfig from path. Returns nil if file missing or invalid.
func LoadAIConfig(path string) (*AIConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var c AIConfig
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// SaveAIConfig writes AIConfig to path, creating parent dirs.
func SaveAIConfig(path string, c *AIConfig) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// ResolveAIConfig returns merged AI config: explicit path > local > global > nil.
// If a path is provided (e.g. from --config flag), only that file is used.
func ResolveAIConfig(explicitPath, cwd string) (*AIConfig, string) {
	if explicitPath != "" {
		c, _ := LoadAIConfig(explicitPath)
		return c, explicitPath
	}
	// Local first
	local := LocalPath(cwd)
	if c, _ := LoadAIConfig(local); c != nil {
		return c, local
	}
	// Then global
	global := GlobalPath()
	if c, _ := LoadAIConfig(global); c != nil {
		return c, global
	}
	return nil, ""
}
