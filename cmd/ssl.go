package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var sslCmd = &cobra.Command{
	Use:   "ssl",
	Short: "Manage SSL certificates",
}

var sslRetrieveCmd = &cobra.Command{
	Use:   "retrieve",
	Short: "Retrieve SSL certificate bundle for a domain",
	Example: `  shadowfax ssl retrieve --domain example.com
  shadowfax ssl retrieve --domain example.com --output /path/to/output`,
	RunE: func(cmd *cobra.Command, args []string) error {
		domain, _ := cmd.Flags().GetString("domain")
		output, _ := cmd.Flags().GetString("output")

		if domain == "" {
			return fmt.Errorf("--domain is required")
		}

		bundle, err := client.RetrieveSSL(domain)
		if err != nil {
			return err
		}

		if outputJSON(cmd) {
			return printJSON(bundle)
		}

		if output != "" {
			dir := output
			if err := os.MkdirAll(dir, 0700); err != nil {
				return fmt.Errorf("creating output directory: %w", err)
			}

			files := map[string]string{
				"domain.cert.pem": bundle.CertificateChain,
				"private.key.pem": bundle.PrivateKey,
				"public.key.pem":  bundle.PublicKey,
			}

			for name, content := range files {
				path := filepath.Join(dir, name)
				if err := os.WriteFile(path, []byte(content), 0600); err != nil {
					return fmt.Errorf("writing %s: %w", name, err)
				}
				if !isQuiet(cmd) {
					fmt.Printf("✓ Wrote %s\n", path)
				}
			}

			return nil
		}

		if isQuiet(cmd) {
			return nil
		}

		// Print to stdout
		fmt.Println("=== Certificate Chain ===")
		fmt.Println(bundle.CertificateChain)
		fmt.Println("=== Private Key ===")
		fmt.Println(bundle.PrivateKey)
		fmt.Println("=== Public Key ===")
		fmt.Println(bundle.PublicKey)

		return nil
	},
}

func init() {
	sslRetrieveCmd.Flags().String("domain", "", "Domain name (e.g. example.com)")
	sslRetrieveCmd.Flags().String("output", "", "Directory to write certificate files to")

	sslCmd.AddCommand(sslRetrieveCmd)
	rootCmd.AddCommand(sslCmd)
}
