package utilities

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:108.0) Gecko/20100101 Firefox/108.0"
)

// Colors for terminal output
var (
	Gray  = "\033[90m"
	Reset = "\033[0m"
)

// Httpx runs the httpx command on a list of subdomains and returns parsed JSON responses
func Httpx(subdomains []string, domain string) ([]map[string]interface{}, error) {
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

	fmt.Printf("%sExecuting command: %s%s\n", Gray, command, Reset)

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
