package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newListRecordsCmd() *cobra.Command {
	var filter string

	cmd := &cobra.Command{
		Use:   "list-records",
		Short: "List DNS records",
		RunE: func(cmd *cobra.Command, args []string) error {
			zone, err := fetchRecords()
			if err != nil {
				return err
			}

			filterRecords(zone, filter)
			return nil
		},
	}

	cmd.Flags().StringVarP(&filter, "filter", "f", "", "filter records by name (regexp)")
	return cmd
}

func newDeleteRecordsCmd() *cobra.Command {
	var filter string

	cmd := &cobra.Command{
		Use:   "delete-records",
		Short: "Delete DNS records",
		RunE: func(cmd *cobra.Command, args []string) error {
			zone, err := fetchRecords()
			if err != nil {
				return err
			}

			recordsToDelete := filterRecordsToDelete(zone, filter)
			if len(recordsToDelete) == 0 {
				fmt.Println("No records found matching the filter")
				return nil
			}

			fmt.Println("The following records will be deleted:")
			for _, record := range recordsToDelete {
				fmt.Printf("%s (%s)\n", record.Name, record.Type)
			}

			force, _ := cmd.Flags().GetBool("force")
			if !force {
				fmt.Printf("About to delete %d records. Are you sure? (y/N): ",
					len(recordsToDelete))
				var response string
				fmt.Scanln(&response)

				if strings.ToLower(response) != "y" {
					fmt.Println("Operation cancelled")
					return nil
				}
			}

			return deleteRecords(recordsToDelete, force)
		},
	}

	cmd.Flags().StringVarP(&filter, "filter", "f", "", "filter records by name (regexp)")
	cmd.Flags().StringVarP(&config.ExcludePattern, "exclude", "e", "", "Exclude records matching this pattern (regexp)")
	cmd.Flags().Bool("force", false, "Skip confirmation prompt when deleting records")
	return cmd
}
