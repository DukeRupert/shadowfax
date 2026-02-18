package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

const defaultTTL = "600"

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Verify credentials and connectivity",
	RunE: func(cmd *cobra.Command, args []string) error {
		ip, err := client.Ping()
		if err != nil {
			return err
		}
		fmt.Printf("✓ Credentials valid. Your IP: %s\n", ip)
		return nil
	},
}

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Manage Porkbun DNS records",
}

var dnsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a DNS record",
	Example: `  shadowfax dns create --domain example.com --type A --content 1.2.3.4
  shadowfax dns create --domain example.com --type CNAME --name www --content example.com
  shadowfax dns create --domain example.com --type A --content 1.2.3.4 --ttl 300`,
	RunE: func(cmd *cobra.Command, args []string) error {
		domain, _ := cmd.Flags().GetString("domain")
		name, _ := cmd.Flags().GetString("name")
		recordType, _ := cmd.Flags().GetString("type")
		content, _ := cmd.Flags().GetString("content")
		ttl, _ := cmd.Flags().GetString("ttl")

		if domain == "" || recordType == "" || content == "" {
			return fmt.Errorf("--domain, --type, and --content are required")
		}

		id, err := client.CreateRecord(domain, name, recordType, content, ttl)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Created %s record for %s", recordType, domain)
		if name != "" {
			fmt.Printf(" (name: %s)", name)
		}
		fmt.Printf(" → ID: %s\n", id)
		return nil
	},
}

var dnsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a DNS record by ID",
	Example: `  shadowfax dns delete --domain example.com --id 123456`,
	RunE: func(cmd *cobra.Command, args []string) error {
		domain, _ := cmd.Flags().GetString("domain")
		id, _ := cmd.Flags().GetString("id")

		if domain == "" || id == "" {
			return fmt.Errorf("--domain and --id are required")
		}

		if err := client.DeleteRecord(domain, id); err != nil {
			return err
		}

		fmt.Printf("✓ Deleted record %s from %s\n", id, domain)
		return nil
	},
}

var dnsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all DNS records for a domain",
	Example: `  shadowfax dns list --domain example.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		domain, _ := cmd.Flags().GetString("domain")

		if domain == "" {
			return fmt.Errorf("--domain is required")
		}

		records, err := client.ListRecords(domain)
		if err != nil {
			return err
		}

		if len(records) == 0 {
			fmt.Printf("No records found for %s\n", domain)
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tTYPE\tNAME\tCONTENT\tTTL")
		fmt.Fprintln(w, "──\t────\t────\t───────\t───")
		for _, r := range records {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", r.ID, r.Type, r.Name, r.Content, r.TTL)
		}
		w.Flush()

		return nil
	},
}

func init() {
	// create flags
	dnsCreateCmd.Flags().String("domain", "", "Domain name (e.g. example.com)")
	dnsCreateCmd.Flags().String("name", "", "Subdomain name (e.g. www) — omit for root record")
	dnsCreateCmd.Flags().String("type", "", "Record type (A, CNAME, MX, TXT, etc.)")
	dnsCreateCmd.Flags().String("content", "", "Record content (e.g. IP address or target domain)")
	dnsCreateCmd.Flags().String("ttl", defaultTTL, "Time to live in seconds (default 600)")

	// delete flags
	dnsDeleteCmd.Flags().String("domain", "", "Domain name (e.g. example.com)")
	dnsDeleteCmd.Flags().String("id", "", "Record ID to delete (get from dns list)")

	// list flags
	dnsListCmd.Flags().String("domain", "", "Domain name (e.g. example.com)")

	// wire up
	dnsCmd.AddCommand(dnsCreateCmd)
	dnsCmd.AddCommand(dnsDeleteCmd)
	dnsCmd.AddCommand(dnsListCmd)
	rootCmd.AddCommand(dnsCmd)
	rootCmd.AddCommand(pingCmd)
}