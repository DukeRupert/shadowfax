package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dukerupert/shadowfax/pkg/porkbun"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("output", "", "Output format: json for machine-readable JSON")
	rootCmd.PersistentFlags().Bool("quiet", false, "Suppress output except errors")
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(os.Getenv("HOME"), ".config", "shadowfax"))

	// Map env vars to viper keys
	viper.SetEnvPrefix("")
	viper.BindEnv("api_key", "PORKBUN_API_KEY")
	viper.BindEnv("secret_key", "PORKBUN_SECRET_KEY")

	// Read config file (ignore error if not found)
	_ = viper.ReadInConfig()
}

func outputJSON(cmd *cobra.Command) bool {
	val, _ := cmd.Flags().GetString("output")
	return val == "json"
}

func isQuiet(cmd *cobra.Command) bool {
	val, _ := cmd.Flags().GetBool("quiet")
	return val
}

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func initClient() error {
	// Load .env files as fallback (lowest priority)
	dotfilePath := filepath.Join(os.Getenv("HOME"), ".dotfiles", ".env")
	if _, err := os.Stat(dotfilePath); err == nil {
		_ = godotenv.Load(dotfilePath)
	} else {
		_ = godotenv.Load(".env")
	}

	// Viper checks: env vars > config file > .env (already loaded into env)
	apiKey := viper.GetString("api_key")
	secretKey := viper.GetString("secret_key")

	// Fall back to env vars directly (covers .env loading)
	if apiKey == "" {
		apiKey = os.Getenv("PORKBUN_API_KEY")
	}
	if secretKey == "" {
		secretKey = os.Getenv("PORKBUN_SECRET_KEY")
	}

	if apiKey == "" || secretKey == "" {
		return fmt.Errorf("API credentials required. Set via:\n  • ~/.config/shadowfax/config.yaml (api_key / secret_key)\n  • Environment variables (PORKBUN_API_KEY / PORKBUN_SECRET_KEY)\n  • ~/.dotfiles/.env or ./.env")
	}

	client = porkbun.NewClient(apiKey, secretKey)
	return nil
}
