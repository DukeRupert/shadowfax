package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `Generate a completion script for the specified shell.

To load completions:

Bash:
  $ source <(shadowfax completion bash)
  # To load completions for each session, execute once:
  $ shadowfax completion bash > /etc/bash_completion.d/shadowfax

Zsh:
  $ source <(shadowfax completion zsh)
  # To load completions for each session, execute once:
  $ shadowfax completion zsh > "${fpath[1]}/_shadowfax"

Fish:
  $ shadowfax completion fish | source
  # To load completions for each session, execute once:
  $ shadowfax completion fish > ~/.config/fish/completions/shadowfax.fish

PowerShell:
  PS> shadowfax completion powershell | Out-String | Invoke-Expression
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
