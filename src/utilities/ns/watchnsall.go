package ns

import (
	"fmt"
	"log"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
)

func NsAllProgram() {
	// Get database connection and the deferred CloseDb function
	db, closeDb, err := database.GetDbAfterInit()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer closeDb() // Ensure that the connection will be closed when the function exits

	programs, err := model.GetAllPrograms(db)
	if err != nil {
		log.Fatalf("Failed to fetch programs: %v", err)
	}

	// Iterate through programs and print their scopes
	for _, program := range programs {
		fmt.Printf("Program: %s\n", program.ProgramName)
		if len(program.Scopes) > 0 {
			fmt.Printf("Scopes:%s \n", program.Scopes)
			for _, scope := range program.Scopes {
				// run ns scope
				Ns(scope)
			}
		} else {
			fmt.Println("No scopes found for this program.")
		}
	}
}
