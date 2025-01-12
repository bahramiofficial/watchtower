package ns

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
)


func Ns(domain string) {
	// Get database connection and the deferred CloseDb function
	db, closeDb, err := database.GetDbAfterInit()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer closeDb() // Ensure that the connection will be closed when the function exits

	subdomain, err := model.GetAllSubdomainWithScopeName(db, domain)
	if err != nil {
		fmt.Printf("can't get subdomain")
	}
	out, err := RunNsCommand(subdomain, domain)
	if err != nil {
		fmt.Printf("can't run ns command")
	}
	fmt.Printf("%v", out)
}

// RunNsCommand executes the dnsx command using a list of subdomains and a domain.
func RunNsCommand(subdomains []string, domain string) (string, error) {
	tmpFile, err := os.CreateTemp("", "subdomains_*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tmpFile.Close()

	// Step 2: Write subdomains to the temporary file
	for _, subdomain := range subdomains {
		if _, err := tmpFile.WriteString(subdomain + "\n"); err != nil {
			return "", fmt.Errorf("failed to write to temporary file: %w", err)
		}
	}

	// Step 3: Get the path of the temporary file
	tmpFilePath := tmpFile.Name()

	// Step 4: Build the dnsx command
	cmd := exec.Command("dnsx", "-l", tmpFilePath, "-silent", "-wd", domain, "-resp", "-json", "-rl", "30", "-t", "10", "-r", "8.8.4.4,129.250.35.251,208.67.222.222")
	cmdOut, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("failed to execute dnsx command: %w\nOutput: %s", err, string(cmdOut))
	}

	// Step 5: Remove the temporary file after running the command
	if err := os.Remove(tmpFilePath); err != nil {
		return "", fmt.Errorf("failed to remove temporary file: %w", err)
	}

	// Step 6: Return the result as a string (this is the dnsx output)
	return string(cmdOut), nil
}
