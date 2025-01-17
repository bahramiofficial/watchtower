package http

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
	"github.com/bahramiofficial/watchtower/src/utilities"
	"gorm.io/gorm"
)

func convertToMapField(input interface{}) model.MapField {
	headers := make(model.MapField)
	for key, value := range input.(map[string]interface{}) {
		headers[key] = fmt.Sprintf("%v", value) // Convert value to string if necessary
	}
	return headers
}

// processResults converts input and calls upsertHTTP
func processResults(db *gorm.DB, results []map[string]interface{}, domain string) {
	for _, obj := range results {
		http := model.Http{
			SubDomain:  obj["input"].(string),
			Scope:      domain,
			IPs:        obj["a"].([]string),
			Tech:       obj["tech"].([]string),
			Title:      obj["title"].(string),
			StatusCode: obj["status_code"].(string),
			Headers:    convertToMapField(obj["header"]),
			URL:        obj["url"].(string),
			FinalURL:   obj["final_url"].(string),
			Favicon:    obj["favicon"].(string),
		}
		model.UpsertHttp(db, http)
	}
}

// 28 5
// if use RunHttpx    add if info cdn of public or private for performance
func Httpx(domain string) {
	// Get database connection and the deferred CloseDb function
	db, closeDb, err := database.GetDbAfterInit()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer closeDb() // Ensure that the connection will be closed when the function exits

	// Get program by scope
	liveSubdomain, err := model.GetAllLiveSubdomainWithScopeName(db, domain)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	res, err := RunHttpx(liveSubdomain, domain)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	processResults(db, res, domain)

} // Httpx runs the httpx command on a list of subdomains and returns parsed JSON responses

func RunHttpx(subdomains []string, domain string) ([]map[string]interface{}, error) {
	userAgent := utilities.GetUserAgent()

	// Create a temporary file to store subdomains
	tempFile, err := os.CreateTemp("", "subdomains_*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name()) // Ensure file is deleted after function returns

	// Write subdomains to the temporary file
	writer := bufio.NewWriter(tempFile)
	for _, sub := range subdomains {
		if _, err := writer.WriteString(sub + "\n"); err != nil {
			return nil, fmt.Errorf("failed to write subdomains to temp file: %w", err)
		}
	}
	writer.Flush()
	tempFile.Close()

	// Construct the command
	command := fmt.Sprintf(
		"httpx -l %s -silent -json -favicon -fhr -tech-detect -irh -include-chain -timeout 3 -retries 1 -threads 5 -rate-limit 4 -ports 443 -extract-fqdn -H 'User-Agent: %s' -H 'Referer: https://%s'",
		tempFile.Name(), userAgent, domain,
	)

	fmt.Printf("Executing command: %s\n", command)

	// Run the command
	cmd := exec.Command("zsh", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute httpx: %w", err)
	}

	// Parse each line of the output as JSON
	var responses []map[string]interface{}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		var response map[string]interface{}
		if err := json.Unmarshal([]byte(scanner.Text()), &response); err != nil {
			return nil, fmt.Errorf("failed to parse JSON response: %w", err)
		}
		responses = append(responses, response)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading command output: %w", err)
	}

	return responses, nil
}
