package ns

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
	"github.com/bahramiofficial/watchtower/src/utilities"
)

// Ns performs DNS queries for subdomains of the specified domain and processes the results.
func Ns(domain string) {
	// Initialize the database connection and ensure it is closed after the function exits
	db, closeDb, err := database.GetDbAfterInit()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer closeDb()

	// Retrieve all subdomains associated with the given domain from the database
	subdomain, err := model.GetAllSubdomainWithScopeName(db, domain)
	if err != nil {
		log.Printf("Can't get subdomain: %v", err)
		return
	}

	fmt.Print(subdomain)
	// Execute the DNS resolution command and handle errors
	results, err := RunNsCommand(subdomain, domain)
	if err != nil {
		log.Printf("Can't run ns command: %v", err)
		return
	}

	// Process each result, parsing JSON data and extracting relevant information
	for _, result := range results {
		var data map[string]interface{} // Dynamic map to hold JSON fields
		if err := json.Unmarshal([]byte(result), &data); err != nil {
			log.Printf("Error parsing JSON: %v", err)
			continue
		}

		// Extract the "a" records from the parsed data and convert them to a string slice
		aRecords, ok := data["a"].([]interface{})
		if !ok {
			log.Printf("Unexpected type for 'a' records: %v", data["a"])
			continue
		}

		ips := make([]string, len(aRecords))
		for i, ip := range aRecords {
			ips[i] = fmt.Sprint(ip)
		}

		// Determine the tag (cdn, public, or private) based on IP classification
		tag, err := utilities.GetIPTag(ips)
		if err != nil {
			log.Printf("Error determining IP tag: %v", err)
			continue
		}

		// Display the extracted host, "a" records, and tag
		var sub = data["host"].(string)

		err = model.UpsertLiveSubdomain(db, domain, sub, ips, tag)
		if err != nil {
			fmt.Printf("not upsert live subdomain  %s%s", domain, sub)
		}
	}

	// Print a separator for visual clarity
	fmt.Println(strings.Repeat("-", 50))
}

// RunNsCommand runs the dnsx command with a list of subdomains and returns the results.
func RunNsCommand(subdomains []string, domain string) ([]string, error) {
	// Create a temporary file to store the list of subdomains
	tmpFile, err := os.CreateTemp("", "subdomains_*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Ensure the file is removed after use
	defer tmpFile.Close()

	// Write each subdomain to the temporary file
	for _, subdomain := range subdomains {
		if _, err := tmpFile.WriteString(subdomain + "\n"); err != nil {
			return nil, fmt.Errorf("failed to write to temporary file: %w", err)
		}
	}

	// Execute the dnsx command with the provided arguments
	cmd := exec.Command("dnsx", "-l", tmpFile.Name(), "-silent", "-wd", domain, "-resp", "-json", "-rl", "30", "-t", "10", "-r", "8.8.4.4,129.250.35.251,208.67.222.222")
	cmdOut, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute dnsx command: %w\nOutput: %s", err, string(cmdOut))
	}

	// Parse the command output line by line into a slice of strings
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOut)))
	var results []string
	for scanner.Scan() {
		results = append(results, scanner.Text())
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return nil, fmt.Errorf("error scanning command output: %w", scanErr)
	}

	return results, nil
}
