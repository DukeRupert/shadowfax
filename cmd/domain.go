package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Manage Porkbun domains",
}

var domainListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all domains in your account",
	Example: `  shadowfax domain list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		domains, err := client.ListDomains()
		if err != nil {
			return err
		}

		if outputJSON(cmd) {
			return printJSON(domains)
		}
		if isQuiet(cmd) {
			return nil
		}

		if len(domains) == 0 {
			fmt.Println("No domains found in your account")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "DOMAIN\tSTATUS\tAUTO-RENEW\tEXPIRES")
		fmt.Fprintln(w, "──────\t──────\t──────────\t───────")
		for _, d := range domains {
			fmt.Fprintf(w, "%s\t%s\t%v\t%s\n", d.Domain, d.Status, d.AutoRenew, d.ExpireDate)
		}
		w.Flush()

		return nil
	},
}

var domainPricingCmd = &cobra.Command{
	Use:   "pricing",
	Short: "Get pricing for a TLD",
	Example: `  shadowfax domain pricing --tld com
  shadowfax domain pricing --tld llc`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tld, _ := cmd.Flags().GetString("tld")

		pricing, err := client.GetPricing()
		if err != nil {
			return err
		}

		if tld != "" {
			p, ok := pricing[tld]
			if !ok {
				return fmt.Errorf("no pricing found for TLD: %s", tld)
			}

			if outputJSON(cmd) {
				return printJSON(map[string]any{tld: p})
			}
			if isQuiet(cmd) {
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "TLD\tREGISTRATION\tRENEWAL\tTRANSFER")
			fmt.Fprintln(w, "───\t────────────\t───────\t────────")
			fmt.Fprintf(w, "%s\t$%s\t$%s\t$%s\n", tld, p.Registration, p.Renewal, p.Transfer)
			w.Flush()
			return nil
		}

		if outputJSON(cmd) {
			return printJSON(pricing)
		}
		if isQuiet(cmd) {
			return nil
		}

		// Show all TLDs sorted
		tlds := make([]string, 0, len(pricing))
		for t := range pricing {
			tlds = append(tlds, t)
		}
		sort.Strings(tlds)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "TLD\tREGISTRATION\tRENEWAL\tTRANSFER")
		fmt.Fprintln(w, "───\t────────────\t───────\t────────")
		for _, t := range tlds {
			p := pricing[t]
			fmt.Fprintf(w, "%s\t$%s\t$%s\t$%s\n", t, p.Registration, p.Renewal, p.Transfer)
		}
		w.Flush()

		return nil
	},
}

func init() {
	// pricing flags
	domainPricingCmd.Flags().String("tld", "", "TLD to get pricing for (e.g. com, llc) — omit for all TLDs")

	// wire up
	domainCmd.AddCommand(domainListCmd)
	domainCmd.AddCommand(domainPricingCmd)
	rootCmd.AddCommand(domainCmd)
}
