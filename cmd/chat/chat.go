package chat

import (
	"context"
	"fmt"
	"os"

	"agentai-go/internal/config"
	"agentai-go/internal/core"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// NewCommand returns the chat subcommand.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chat [goal]",
		Short: "Generate code based on a user goal",
		Long:  "Agentic AI Code Assistant - provide a goal and the tool will plan and generate code.",
		Args:  cobra.ExactArgs(1),
		RunE:  runChat,
	}
	return cmd
}

func runChat(cmd *cobra.Command, args []string) error {
	goal := args[0]
	fmt.Printf("\n🎯 Processing goal: %s\n\n", goal)

	_ = godotenv.Load() // .env in cwd
	configPath, _ := cmd.Flags().GetString("config")
	cwd, _ := os.Getwd()
	cfg := config.Load(configPath, cwd)
	apiKey := cfg.APIKey
	if apiKey == "" {
		apiKey = cfg.GeminiAPIToken
	}
	if apiKey == "" || apiKey == "YOUR_GEMINI_API_TOKEN_HERE" {
		fmt.Fprintln(os.Stderr, "Error: No API key. Set with 'agentai config set api_key <key>' (--local or --global) or .env (GEMINI_API_TOKEN / OPENAI_API_KEY).")
		os.Exit(1)
	}

	ctx := context.Background()
	if err := core.ChatCommand(ctx, goal, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "\n❌ Error: %v\n", err)
		os.Exit(1)
	}
	return nil
}
