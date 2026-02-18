package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dukerupert/shadowfax/pkg/porkbun"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var client *porkbun.Client

var rootCmd = &cobra.Command{
	Use:   "shadowfax",
	Short: "A CLI tool for managing Porkbun DNS records",
	Long:  `Shadowfax manages Porkbun DNS records via the Porkbun JSON API.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initClient()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initClient() error {
	// Try ~/.dotfiles/.env first, fall back to .env in current directory
	dotfilePath := filepath.Join(os.Getenv("HOME"), ".dotfiles", ".env")
	if _, err := os.Stat(dotfilePath); err == nil {
		_ = godotenv.Load(dotfilePath)
	} else {
		_ = godotenv.Load(".env")
	}

	apiKey := os.Getenv("PORKBUN_API_KEY")
	secretKey := os.Getenv("PORKBUN_SECRET_KEY")

	if apiKey == "" || secretKey == "" {
		return fmt.Errorf("PORKBUN_API_KEY and PORKBUN_SECRET_KEY must be set in ~/.dotfiles/.env or .env")
	}

	client = porkbun.NewClient(apiKey, secretKey)
	return nil
}