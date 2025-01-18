package watchhttp

import (
	"fmt"
	"log"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
)

func HttpAll() {
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
	// http all Program scopes
	for _, program := range programs {
		fmt.Printf("http all Program \n")
		for _, domain := range program.Scopes {
			Httpx(domain)
		}
	}
}
