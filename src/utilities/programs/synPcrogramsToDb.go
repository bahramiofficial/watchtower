package watchsyncprograms

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
	// Import the required package
)

// for call to api app
func SyncProgramToDb(dirPath string) {

	// Initialize the database connectionn
	err := database.InitDb()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDb()

	db := database.GetDb()

	// dynamic path
	if dirPath == "" {
		// Get the directory of the script
		_, scriptPath, _, ok := runtime.Caller(0)
		if !ok {
			log.Fatalf("Failed to get the script directory")
		}

		// init path json files
		dirPath = filepath.Join(filepath.Dir(scriptPath), "json")
	}

	// Is Exist directory
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		log.Fatalf("Directory %s does not exist", dirPath)
	}

	// Scan the directory
	files, err := filepath.Glob(filepath.Join(dirPath, "*.json"))
	if err != nil {
		log.Fatalf("Failed to scan directory: %v", err)
	}

	if len(files) == 0 {
		log.Println("No JSON files found.")
		return
	}

	// Process each JSON file
	for _, file := range files {
		fmt.Printf("Processing file: %s\n", file)

		// Read the file content
		content, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Failed to read file %s: %v\n", file, err)
			continue
		}

		// Parse the JSON content
		var program model.Program
		if err := json.Unmarshal(content, &program); err != nil {
			log.Printf("Failed to parse JSON in file %s: %v\n", file, err)
			continue
		}

		// Insert or update the database record
		err = model.InsertOrUpdateProgram(db, &program)
		if err != nil {
			log.Printf("Failed to insert or update record for file %s: %v\n", file, err)
		} else {
			fmt.Printf("Successfully inserted or updated record for program: %s\n", program.ProgramName)
		}

	}
}
