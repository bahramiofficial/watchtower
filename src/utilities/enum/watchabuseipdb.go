package watch

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
	"github.com/bahramiofficial/watchtower/src/utilities"
)

func RunAbuseIPDB(domain string) {

	// Get database connection and the deferred CloseDb function
	db, closeDb, err := database.GetDbAfterInit()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer closeDb() // Ensure that the connection will be closed when the function exits

	// Check if domain is empty
	if domain == "" {
		fmt.Println("Usage: watch_abuseipdb domain")
		return
	}

	// Get program by scope
	program, err := model.GetProgramByScope(db, domain)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if program != nil {
		// Log the action
		fmt.Printf("[%s] Running AbuseIPDB module for '%s'\n", utilities.GetFormattedTime(), domain)

		// Call abuseipdb function to get subdomains
		subs := AbuseIPDB(domain)

		// Loop through the subdomains and upsert them if they match the domain
		for _, sub := range subs {
			if sub != "" {
				// Upsert subdomain if match found
				err := model.UpsertSubdomain(db, program.ProgramName, sub, "abuseipdb")
				if err != nil {
					log.Printf("Failed to upsert subdomain '%s': %v", sub, err)
				}
			}
		}
	} else {
		// Log that the scope for the domain doesn't exist
		fmt.Printf("[%s] Scope for '%s' does not exist in Watchtower\n", utilities.GetFormattedTime(), domain)
	}
}

// AbuseIPDB queries abuseipdb for a given domain and processes the results.
func AbuseIPDB(domain string) []string {
	// Define the command with the appropriate user-agent and session token
	command := fmt.Sprintf(
		`curl -s "https://www.abuseipdb.com/whois/%s" -H "user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36" -b "abuseipdb_session=YOUR-SESSION" | grep --color=auto --exclude-dir={.bzr,CVS,.git,.hg,.svn,.idea,.tox} -E "<li>\\w.*</li>" | sed -E "s/<\\/?li>//g" | sed "s|$|.%s|"`,
		domain, domain,
	)

	fmt.Printf("%s Executing command: %s\n", utilities.GetFormattedTime(), command)

	// Execute the command
	cmd := exec.Command("zsh", "-c", command)
	output, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get output pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start command: %v", err)
	}

	// Process the command output
	var results []string
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line != "" {
			results = append(results, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading command output: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Command execution error: %v", err)
	}

	resNum := len(results)
	fmt.Printf("%s done for %s, results: %d\n", utilities.GetFormattedTime(), domain, resNum)
	return results
}
