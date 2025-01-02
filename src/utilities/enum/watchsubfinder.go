package watch

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
	"github.com/bahramiofficial/watchtower/src/utilities"
	// Import SQLite3 driver
)

func Subfinder(domain string, outputFile string) {
	if outputFile == "" {
		outputFile = "/Users/mrpit/Documents/GitHub/watchtower/src/cmd/subdomains.txt"
	}

	err := database.InitDb()
	if err != nil {
		log.Fatalf("%v Failed to initialize database: %v", utilities.GetFormattedTime(), err)

	}
	defer database.CloseDb()

	db := database.GetDb()
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
