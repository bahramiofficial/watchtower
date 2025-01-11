package watch

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
	"github.com/bahramiofficial/watchtower/src/utilities"
)

func Chaos(domain string) {
	// Get database connection and the deferred CloseDb function
	db, closeDb, err := database.GetDbAfterInit()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer closeDb() // Ensure that the connection will be closed when the function exits

	// Fetch the program by scope
	program, err := model.GetProgramByScope(db, domain)

	if err == nil {
		// Print the current time and start Chaos module
		fmt.Printf("[%s] Running Chaos module for '%s'\n", utilities.GetFormattedTime(), domain)

		// Fetch subdomains from the Chaos folder
		subs := Runchaos(domain)
		for _, sub := range subs {
			if err == nil {
				model.UpsertSubdomain(db, program.ProgramName, sub, "chaos")
			}
		}
	} else {
		// Log if the program for the given domain does not exist
		fmt.Printf("[%s] Scope for '%s' does not exist in Watchtower\n", utilities.GetFormattedTime(), domain)
	}

	// Download and unzip Chaos data
	downloadAndUnzip()
}

// Chaos function reads the domain-specific file in the CHAOS_FOLDER
func Runchaos(domain string) []string {
	// Construct the path to the file in the CHAOS_FOLDER
	filePath := filepath.Join(utilities.GetRootPath(), "/data/chaos/", fmt.Sprintf("%s.txt", domain))

	// Check if the file exists
	if _, err := os.Stat(filePath); err == nil {
		// Open the file if it exists
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return nil
		}
		defer file.Close()

		var lines []string
		// Read lines from the file
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				lines = append(lines, line)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file: %v\n", err)
		}

		return lines
	} else {
		// If the file does not exist
		fmt.Printf("%s does not exist in Chaos module.\n", domain)
		return nil
	}
}

// Download and unzip files from the Chaos project
func downloadAndUnzip() {
	// Execute the command to download and unzip Chaos data
	curlCommand := fmt.Sprintf(
		"curl -s https://chaos-data.projectdiscovery.io/index.json | " +
			"jq -r \".[].URL\" | " +
			"while read url; do " +
			"wget -q $url && " +
			"unzip -q $(echo $url | rev | cut -d / -f 1 | rev); done && " +
			"rm -rf *.zip",
	)

	// Run the command using zsh
	cmd := exec.Command("zsh", "-c", curlCommand)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		return
	}
	fmt.Println(string(output))

	// Move the downloaded files to CHAOS_FOLDER
	files, err := os.ReadDir("./")
	if err != nil {
		fmt.Printf("Error reading temp directory: %v\n", err)
		return
	}
	CHAOS_FOLDER := filepath.Join(utilities.GetRootPath(), "/data/chaos")

	for _, file := range files {
		if !file.IsDir() {
			src := file.Name()
			dst := filepath.Join(CHAOS_FOLDER, file.Name())
			err := os.Rename(src, dst)
			if err != nil {
				fmt.Printf("Error moving file: %v\n", err)
			}
		}
	}

	fmt.Println("Files downloaded and unzipped successfully.")
}
