package watch

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
	"github.com/bahramiofficial/watchtower/src/utilities"
	// Import SQLite3 driver
)

func Subfinder(domain string, outputFile string) {
	// Get database connection and the deferred CloseDb function
	db, closeDb, err := database.GetDbAfterInit()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer closeDb() // Ensure that the connection will be closed when the function exits

	if outputFile == "" {
		outputFile = filepath.Join(utilities.GetRootDirPath(), "data", "subfinder", "subfindersubs.txt")
	}

	// Check if file exists
	if _, err := os.Stat(outputFile); err == nil {
		// File exists, remove it
		err := os.Remove(outputFile)
		if err != nil {
			fmt.Print("not removing", outputFile)
		}
	}

	var program *model.Program
	var results *[]string

	// Define the domain
	program, _ = model.FindDomainWithProgramName(db, domain)
	if program != nil {
		fmt.Printf("%v running subfinder on %v\n", utilities.GetFormattedTime(), domain)
		results, err = RunSubfinder(domain, outputFile)
		if err != nil {
			fmt.Printf("Failed to run Subfinder: %v", err)
		}
		// Check if there are any results
		if results == nil || len(*results) == 0 {
			fmt.Println("No results found.")
			return
		}

		for _, subdomain := range *results {
			if subdomain != "" {
				model.UpsertSubdomain(db, program.ProgramName, subdomain, "subfinder")
			}
		}
	} else {

	}

}

func RunSubfinder(domain string, outputFile string) (*[]string, error) {
	cmd := exec.Command("subfinder", "-d", domain, "-o", outputFile, "-silent", "-all")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v output file address: %s\n", utilities.GetFormattedTime(), output)
	results := strings.Split(string(output), "\n")
	return &results, nil
}
