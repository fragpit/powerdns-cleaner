package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func buildURL(base, path string) string {
	base = strings.TrimRight(base, "/")
	path = strings.TrimLeft(path, "/")
	return fmt.Sprintf("%s/%s", base, path)
}

func fetchRecords() (*PowerDNSZone, error) {
	url := buildURL(config.APIURL, fmt.Sprintf("servers/localhost/zones/%s", config.Zone))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", config.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d",
			resp.StatusCode)
	}

	var zone PowerDNSZone
	if err := json.NewDecoder(resp.Body).Decode(&zone); err != nil {
		return nil, err
	}

	return &zone, nil
}

func shouldExcludeRecord(record RRSet, excludePattern string) bool {
	if excludePattern == "" {
		return false
	}

	match, err := regexp.MatchString(excludePattern, record.Name)
	if err != nil {
		// If pattern is invalid, log warning and don't exclude
		fmt.Printf("Warning: Invalid exclude pattern '%s': %v\n",
			excludePattern, err)
		return false
	}
	return match
}

func filterExcludedRecords(records []RRSet, excludePattern string) []RRSet {
	// Return original records if exclude pattern is empty or not defined
	if excludePattern == "" || config.ExcludePattern == "" {
		return records
	}

	var filtered []RRSet
	for _, record := range records {
		if shouldExcludeRecord(record, excludePattern) {
			fmt.Printf("Excluding record: Name=%s, Type=%s\n",
				record.Name, record.Type)
			continue
		}
		filtered = append(filtered, record)
	}

	fmt.Printf("\nTotal records: %d, After exclusion: %d\n\n",
		len(records), len(filtered))
	return filtered
}

func deleteRecords(records []RRSet, force bool) error {
	if len(records) == 0 {
		fmt.Println("No records to delete after applying exclusion filter")
		return nil
	}

	if !force {
		fmt.Printf("About to delete %d records. Are you sure? (y/N): ",
			len(records))
		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) != "y" {
			fmt.Println("Operation cancelled")
			return nil
		}
	}

	url := buildURL(config.APIURL, fmt.Sprintf("servers/localhost/zones/%s",
		config.Zone))

	type DeleteSet struct {
		Name       string `json:"name"`
		Type       string `json:"type"`
		ChangeType string `json:"changetype"`
	}

	const batchSize = 30
	for i := 0; i < len(records); i += batchSize {
		end := i + batchSize
		if end > len(records) {
			end = len(records)
		}

		var deleteSets struct {
			RRSets []DeleteSet `json:"rrsets"`
		}

		fmt.Printf("\nProcessing batch %d-%d:\n", i, end)
		for _, record := range records[i:end] {
			fmt.Printf("- Deleting record: Name=%s, Type=%s\n",
				record.Name, record.Type)
			deleteSets.RRSets = append(deleteSets.RRSets, DeleteSet{
				Name:       record.Name,
				Type:       record.Type,
				ChangeType: "DELETE",
			})
		}

		jsonData, err := json.Marshal(deleteSets)
		if err != nil {
			return fmt.Errorf("failed to marshal batch %d-%d: %w", i, end, err)
		}

		req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create request for batch %d-%d: %w", i, end, err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", config.APIKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to execute request for batch %d-%d: %w", i, end, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			return fmt.Errorf("failed to delete batch %d-%d: status code %d", i, end, resp.StatusCode)
		}

		fmt.Printf("✓ Successfully deleted batch %d-%d\n", i, end)
	}

	return nil
}
