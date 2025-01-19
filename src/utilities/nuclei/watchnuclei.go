package nuclei

import (
	"fmt"
	"log"
	"os"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
	"github.com/bahramiofficial/watchtower/src/utilities"
)

func Nuclei(domain string) {

	// Get database connection and the deferred CloseDb function
	db, closeDb, err := database.GetDbAfterInit()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer closeDb() // Ensure that the connection will be closed when the function exits
	httpsubdomains, err := model.GetAllHttpWithScope(db, domain)
	if err != nil {
		log.Printf("Error fetching HTTP records with scope '%s': %v", domain, err)
		return
	}

	if len(httpsubdomains) > 0 {
		fmt.Printf("[%s] running Nuclei module for all HTTP services\n", utilities.GetFormattedTime())

		// Collect URLs for Nuclei
		var urls []string
		for _, httpsubdomain := range httpsubdomains {
			urls = append(urls, httpsubdomain.URL) // Adjust to actual field name for URL
		}

		// Execute Nuclei
		results, err := RunNucleiCommand(urls)
		if err != nil {
			fmt.Printf("Error running Nuclei: %v\n", err)
			os.Exit(1)
		}

		// Send results if not empty
		if results != "" {
			//todo
			utilities.SendDiscordMessage(results)
		}
	}
}

func RunNucleiCommand(urls []string) (string, error) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "urls-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name()) // Ensure the temporary file is removed

	// Write URLs to the temporary file
	for _, url := range urls {
		_, err = tempFile.WriteString(fmt.Sprintf("%s\n", url))
		if err != nil {
			return "", fmt.Errorf("failed to write to temporary file: %w", err)
		}
	}
	tempFile.Close()

	// Define command to run nuclei
	rootDirPath := utilities.GetRootDirPath() // Assume RootDirPath() returns the correct path
	command := fmt.Sprintf("nuclei -l %s -config %s/data/nuclei/public-config.yaml", tempFile.Name(), rootDirPath)
	fmt.Printf("Executing command: %s\n", command)

	// Execute command using RunCommandInZsh
	results, err := utilities.RunCommandInZsh(command)
	if err != nil {
		return "", fmt.Errorf("error executing nuclei command: %w", err)
	}

	return results, nil
}
