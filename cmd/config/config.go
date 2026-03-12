package config

import (
	"fmt"
	"os"
	"strings"

	"agentai-go/internal/config"

	"github.com/spf13/cobra"
)

var local, global bool

// NewCommand returns the config subcommand.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage AI provider config (local or global)",
		Long:  "Config is stored in .agentai/config.json. Use --local for repo (./.agentai/config.json) or --global for home (~/.agentai/config.json). Supports: gemini, openai, openrouter, ollama.",
	}
	cmd.PersistentFlags().BoolVar(&local, "local", false, "use repository .agentai/config.json")
	cmd.PersistentFlags().BoolVar(&global, "global", false, "use ~/.agentai/config.json")

	cmd.AddCommand(newShowCommand())
	cmd.AddCommand(newSetCommand())
	return cmd
}

func configPath() (string, error) {
	if local && global {
		return "", fmt.Errorf("use either --local or --global, not both")
	}
	if global {
		return config.GlobalPath(), nil
	}
	// default: local (.agentai/config.json in current directory)
	cwd, _ := os.Getwd()
	return config.LocalPath(cwd), nil
}

func newShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current config",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := configPath()
			if err != nil {
				return err
			}
			scope := "local"
			if global {
				scope = "global"
			}
			fmt.Fprintf(os.Stdout, "Config (%s): %s\n", scope, path)
			c, err := config.LoadAIConfig(path)
			if err != nil {
				return err
			}
			if c == nil {
				fmt.Fprintln(os.Stdout, "(no config file yet)")
				return nil
			}
			key := c.APIKey
			if key != "" {
				if len(key) > 8 {
					key = key[:4] + "..." + key[len(key)-4:]
				} else {
					key = "***"
				}
			}
			fmt.Fprintf(os.Stdout, "  provider: %s\n", c.Provider)
			fmt.Fprintf(os.Stdout, "  api_key:  %s\n", key)
			fmt.Fprintf(os.Stdout, "  model:    %s\n", c.Model)
			if c.BaseURL != "" {
				fmt.Fprintf(os.Stdout, "  base_url: %s\n", c.BaseURL)
			}
			return nil
		},
	}
}

func newSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set a config key (provider, api_key, model, base_url)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := configPath()
			if err != nil {
				return err
			}
			key := strings.ToLower(strings.TrimSpace(args[0]))
			value := strings.TrimSpace(args[1])
			allowed := map[string]bool{"provider": true, "api_key": true, "model": true, "base_url": true}
			if !allowed[key] {
				return fmt.Errorf("invalid key %q (use: provider, api_key, model, base_url)", key)
			}
			c, _ := config.LoadAIConfig(path)
			if c == nil {
				c = &config.AIConfig{Provider: "gemini", Model: "gemini-2.5-flash"}
			}
			switch key {
			case "provider":
				c.Provider = value
			case "api_key":
				c.APIKey = value
			case "model":
				c.Model = value
			case "base_url":
				c.BaseURL = value
			}
			if err := config.SaveAIConfig(path, c); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Updated %s = %s\n", key, maskIfKey(key, value))
			return nil
		},
	}
}

func maskIfKey(key, value string) string {
	if key == "api_key" && len(value) > 8 {
		return value[:4] + "..." + value[len(value)-4:]
	}
	if key == "api_key" {
		return "***"
	}
	return value
}
