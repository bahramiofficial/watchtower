package nuclei

import (
	"fmt"
	"log"
	"os"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
	"github.com/bahramiofficial/watchtower/src/utilities"
)

func NucleiAll() {

	// Get database connection and the deferred CloseDb function
	db, closeDb, err := database.GetDbAfterInit()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer closeDb() // Ensure that the connection will be closed when the function exits

	httpsubdomains, err := model.GetAllHttp(db)
	if err != nil {
		log.Printf("Error fetching all HTTP records: %v", err)
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
