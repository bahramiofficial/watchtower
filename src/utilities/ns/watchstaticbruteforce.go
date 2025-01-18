package ns

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
	"github.com/bahramiofficial/watchtower/src/utilities"
)

func StaticBrute(domain string) {
	// Get database connection and the deferred CloseDb function
	db, closeDb, err := database.GetDbAfterInit()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer closeDb() // Ensure that the connection will be closed when the function exits

	var program *model.Program
	var results *[]string

	// Define the domain
	program, _ = model.FindDomainWithProgramName(db, domain)

	if program != nil {
		fmt.Printf("%v running DynamicBrute on %v\n", utilities.GetFormattedTime(), domain)
		results, err = RunStaticBrute(domain)
		if err != nil {
			fmt.Printf("Failed to run DynamicBrute: %v", err)
		}
		// Check if there are any results
		if results == nil || len(*results) == 0 {
			fmt.Println("No results found.")
			return
		}

		for _, subdomain := range *results {
			if subdomain != "" {
				model.UpsertSubdomain(db, program.ProgramName, subdomain, "staticbrute")
				model.UpsertLiveSubdomain(db, domain, subdomain, nil, "")
				// upsert_lives(domain=domain, subdomain=sub, ips=[], tag="")
			}
		}
	} else {

	}
}
func RunStaticBrute(domain string) (*[]string, error) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "static_brute")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up the temporary directory

	// Define paths for temporary files
	bestDnsPath := filepath.Join(tempDir, "best-dns-wordlist.txt")
	subdomainsPath := filepath.Join(tempDir, "2m-subdomains.txt")
	crunchPath := filepath.Join(tempDir, "4-lower.txt")
	staticWordsPath := filepath.Join(tempDir, "static-finals.txt")
	domainStaticPath := filepath.Join(tempDir, fmt.Sprintf("%s.static", domain))

	// Step 1: Prepare wordlist for static brute
	commands := []string{
		fmt.Sprintf("curl -s https://wordlists-cdn.assetnote.io/data/manual/best-dns-wordlist.txt -o %s", bestDnsPath),
		fmt.Sprintf("curl -s https://wordlists-cdn.assetnote.io/data/manual/2m-subdomains.txt -o %s", subdomainsPath),
		fmt.Sprintf("crunch 1 4 abcdefghijklmnopqrstuvwxyz1234567890 > %s", crunchPath),
		fmt.Sprintf("cat %s %s %s | sort -u > %s", bestDnsPath, subdomainsPath, crunchPath, staticWordsPath),
		fmt.Sprintf("awk -v domain='%s' '{print $0\".\"domain}' %s > %s", domain, staticWordsPath, domainStaticPath),
	}

	for _, cmd := range commands {
		fmt.Printf("Executing command: %s\n", cmd)
		if _, err := utilities.RunCommandInZsh(cmd); err != nil {
			return nil, fmt.Errorf("failed to execute command: %s, error: %w", cmd, err)
		}
	}

	// Step 2: Run shuffledns
	shufflednsCommand := fmt.Sprintf(
		"shuffledns -list %s -d %s -r ~/.resolvers -m $(which massdns) -mode resolve -t 100 -silent",
		domainStaticPath, domain,
	)
	fmt.Printf("Executing command: %s\n", shufflednsCommand)
	output, err := utilities.RunCommandInZsh(shufflednsCommand)
	if err != nil {
		return nil, fmt.Errorf("failed to execute shuffledns command: %w", err)
	}

	results := strings.Split(string(output), "\n")
	return &results, nil

}
