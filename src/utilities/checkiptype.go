package utilities

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// isPrivateIP checks if the provided IP address is a private IPv4 address.
func isPrivateIP(ip string) bool {
	privateIPRegex := regexp.MustCompile(`(?m)^(?:10(?:\.\d{1,3}){3}|172\.(?:1[6-9]|2\d|3[01])(?:\.\d{1,3}){2}|192\.168(?:\.\d{1,3}){2})$`)
	return privateIPRegex.MatchString(ip)
}

// GetIPTag determines whether IPs belong to "cdn", "public", or "private" categories.
func GetIPTag(ips []string) (string, error) {
	// Create a temporary file for storing IPs
	tmpFile, err := os.CreateTemp("", "ips_*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Ensure file is removed

	// Write IPs to temporary file
	for _, ip := range ips {
		if _, err := tmpFile.WriteString(ip + "\n"); err != nil {
			return "", fmt.Errorf("failed to write to temporary file: %w", err)
		}
	}
	tmpFile.Close()

	// Run cut-cdn command
	// cmd := exec.Command("zsh", "-c", fmt.Sprintf("cut-cdn -i %s --silent", tmpFile.Name()))
	cmd := exec.Command("zsh", "-c", fmt.Sprintf("export PATH=$PATH:/Users/mrpit/go/bin && cut-cdn -i %s --silent", tmpFile.Name()))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "cdn", fmt.Errorf("error running cut-cdn: %v\nOutput: %s", err, string(output))
	}

	// Process command output
	results := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(results) != len(ips) {
		return "cdn", nil
	}

	// Determine IP tag based on privacy status
	for _, ip := range ips {
		if !isPrivateIP(ip) {
			return "public", nil
		}
	}
	return "private", nil
}
