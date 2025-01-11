package watch

import (
	"fmt"
	"log"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
)

func EnumAll() {
	// Initialize the database connection
	if err := database.InitDb(); err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}
	defer database.CloseDb() // Ensure the database connection is closed when the application exits

	// Get the database client
	db := database.GetDb()

	// Fetch all programs
	programs, err := model.GetAllPrograms(db)
	if err != nil {
		log.Fatalf("Failed to retrieve programs: %v", err)
	}

	// enum all Program scopes
	for _, program := range programs {
		fmt.Printf("enum all Program : %s\n", program.ProgramName)
		for _, scope := range program.Scopes {
			fmt.Printf("run Subfinder and crtsh on : %s\n", scope)
			Subfinder(scope, "")
			Crtsh(scope)
			RunAbuseIPDB(scope)
			Wayback(scope)
		}
	}

}
