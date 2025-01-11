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

func Wayback(domain string) {
	// Get database connection and the deferred CloseDb function
	db, closeDb, err := database.GetDbAfterInit()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer closeDb() // Ensure that the connection will be closed when the function exits

	// Get program by scope
	program, err := model.GetProgramByScope(db, domain)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if program != nil {
		// Log the action
		fmt.Printf("[%s] Running Wayback module for '%s'\n", utilities.GetFormattedTime(), domain)

		// Call wayback function to get subdomains
		subs, err := RunWaybackURLs(domain)
		if err != nil {
			log.Printf("Error running Wayback: %v", err)
			return
		}

		// Loop through the subdomains and upsert them if they match the domain
		for _, sub := range subs {
			if sub != "" {
				// Upsert subdomain if match found
				err := model.UpsertSubdomain(db, program.ProgramName, sub, "wayback")
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

// RunWaybackURLs runs the shell command equivalent to waybackurls and unfurl for the given domain.
func RunWaybackURLs(domain string) ([]string, error) {
	// Construct the command
	command := fmt.Sprintf("echo %s | waybackurls | unfurl domain | sort -u", domain)

	// Log the command being executed
	fmt.Printf("%s Executing command: %s\n", utilities.GetFormattedTime(), command)

	// Execute the command using zsh
	cmd := exec.Command("zsh", "-c", command)
	output, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get output pipe: %w", err)
	}

	// Start the command execution
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	// Read the output
	var results []string
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			results = append(results, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading command output: %w", err)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("command execution failed: %w", err)
	}

	// Log the number of results
	resNum := len(results)
	fmt.Printf("%s done for %s, results: %d\n", utilities.GetFormattedTime(), domain, resNum)

	return results, nil
}
