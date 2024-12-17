package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var config struct {
	Debug          bool
	APIURL         string
	APIKey         string
	Zone           string
	Filter         string
	ExcludePattern string
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "pdnsutil",
		Short: "PowerDNS CLI utility",
		Long: `pdnsutil is a CLI utility for managing PowerDNS records.
It allows listing and deleting DNS records based on filters.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Read from environment variables first
			if env := os.Getenv("PDNS_API_URL"); env != "" {
				config.APIURL = env
			}
			if env := os.Getenv("PDNS_API_KEY"); env != "" {
				config.APIKey = env
			}
			if env := os.Getenv("PDNS_ZONE"); env != "" {
				config.Zone = env
			}

			// Check if required values are set (either via flags or env vars)
			if config.APIURL == "" {
				return fmt.Errorf("API URL is required (use --api-url flag or PDNS_API_URL env var)")
			}
			if config.APIKey == "" {
				return fmt.Errorf("API key is required (use --api-key flag or PDNS_API_KEY env var)")
			}
			if config.Zone == "" {
				return fmt.Errorf("Zone is required (use --zone flag or PDNS_ZONE env var)")
			}
			return nil
		},
	}

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&config.Debug, "debug", "d", false,
		"enable debug mode")
	rootCmd.PersistentFlags().StringVarP(&config.APIURL, "api-url", "a", "",
		"PowerDNS API URL")
	rootCmd.PersistentFlags().StringVarP(&config.APIKey, "api-key", "k", "",
		"PowerDNS API key")
	rootCmd.PersistentFlags().StringVarP(&config.Zone, "zone", "z", "",
		"zone name")

	// Add subcommands
	rootCmd.AddCommand(newListRecordsCmd())
	rootCmd.AddCommand(newDeleteRecordsCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
