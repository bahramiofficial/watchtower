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

// Gau function retrieves subdomains for a given domain using the "gau" tool.
func Gau(domain string) {
	// Get the database connection
	db, closeDb, err := database.GetDbAfterInit()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer closeDb() // Ensure the connection will be closed when the function exits

	// Ensure the domain is not empty
	if domain == "" {
		log.Printf("Usage: watch_gau domain")
		return
	}

	// Get the program by scope
	program, err := model.GetProgramByScope(db, domain)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if program != nil {
		// Log the action
		log.Printf("[%s] Running Gau module for '%s'\n", utilities.GetFormattedTime(), domain)

		// Call Gau function to get subdomains
		subs, err := RunGau(domain)
		if err != nil {
			log.Printf("Error running Gau: %v", err)
			return
		}

		// Loop through the subdomains and upsert them if they match the domain
		for _, sub := range subs {
			if sub != "" {
				err := model.UpsertSubdomain(db, program.ProgramName, sub, "gau")
				if err != nil {
					log.Printf("Failed to upsert subdomain '%s': %v", sub, err)
				}
			}
		}
	} else {
		// Log that the scope for the domain doesn't exist
		log.Printf("[%s] Scope for '%s' does not exist in Watchtower\n", utilities.GetFormattedTime(), domain)
	}
}

// RunGau runs the 'gau' command and returns a list of subdomains
func RunGau(domain string) ([]string, error) {
	// Construct the command
	command := fmt.Sprintf("gau %s --threads 10 --subs | unfurl domain | sort -u", domain)

	// Log the command being executed
	log.Printf("[%s] Executing command: %s\n", utilities.GetFormattedTime(), command)

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

	// Log the results
	resNum := len(results)
	log.Printf("[%s] Done for '%s', results: %d\n", utilities.GetFormattedTime(), domain, resNum)

	return results, nil
}
