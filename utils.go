package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func filterRecords(zone *PowerDNSZone, filter string) {
	count := 0
	for _, rrset := range zone.RRSets {
		if filter == "" || matchFilter(rrset.Name, filter) {
			for _, record := range rrset.Records {
				fmt.Printf("%s (%s): %s\n", rrset.Name, rrset.Type,
					record.Content)
				count++
			}
		}
	}
	fmt.Printf("\nTotal records: %d\n", count)
}

func filterRecordsToDelete(zone *PowerDNSZone, filter string) []RRSet {
	var recordsToDelete []RRSet
	for _, rrset := range zone.RRSets {
		if filter == "" || matchFilter(rrset.Name, filter) {
			recordsToDelete = append(recordsToDelete, rrset)
		}
	}

	return filterExcludedRecords(recordsToDelete, config.ExcludePattern)
}

func matchFilter(name, filter string) bool {
	if filter == "" {
		return true
	}
	matched, err := regexp.MatchString(filter, name)
	if err != nil {
		return false
	}
	return matched
}

func confirmDeletion(recordCount int) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Are you sure you want to delete these %d records? (yes/no): ", recordCount)
	response, _ := reader.ReadString('\n')
	return strings.ToLower(strings.TrimSpace(response)) == "yes"
}
