package utilities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// ExtractBaseDomain extracts the base domain from a full subdomain.
// It works by splitting the domain and considering the last two segments
// as the base domain.
func ExtractBaseDomain(subdomain string) string {
	// Split the domain by periods (.)
	parts := strings.Split(subdomain, ".")

	// Determine the number of parts in the domain
	numParts := len(parts)

	// If there are less than two parts, return the input as is
	if numParts <= 1 {
		return subdomain
	}

	// If the domain has more than two parts, extract the last two parts as the base domain
	// For example: "x.x.x.x.com" -> "x.com"
	if numParts > 2 {
		return parts[numParts-2] + "." + parts[numParts-1]
	}

	// If it's just a two-part domain, return it as is (e.g., "example.com")
	return subdomain
}

// GetFormattedTime returns the current time in the format "y/m/d : h/m/s"
func GetFormattedTime() string {
	// Get current time
	currentTime := time.Now()

	// Format the current time 2025/01/02 22:18:42
	return currentTime.Format("2006/01/02   15:04:05")
}

func RunCommandInZsh(command string) (string, error) {
	cmd := exec.Command("zsh", "-c", command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func SendDiscordMessage(message string) {
	// Define the payload structure
	payload := map[string]string{
		"content": message,
	}

	// Convert the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("failed to marshal JSON: %v", err)
	}

	// Replace with your function to get configuration values
	webhookURL := GetWebHookDiscordUrl()
	if webhookURL == "" {
		fmt.Printf("webhook URL is not configured")
	} else {
		// Send the POST request
		resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("failed to send POST request: %v", err)
		}
		defer resp.Body.Close()

		// Check the status code
		if resp.StatusCode != 204 {
			fmt.Printf("failed to send message. Status code: %d", resp.StatusCode)
		}
		fmt.Printf("Status code: %d", resp.StatusCode)
	}

}
func GetUserAgent() string {
	return "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:108.0) Gecko/20100101 Firefox/108.0"
}

// env
func GetValueEnv(key string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return ""
	}

	err = godotenv.Load(filepath.Join(currentDir, "config.env"))
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	value := os.Getenv(key)
	return value
}
func GetWebHookDiscordUrl() string {
	return GetValueEnv("WEBHOOK_DISCORD_URL")
}
func GetRootDirPath() string {
	return GetValueEnv("ROOT_PATH_DIR")
}

// AreSlicesEqual compares two slices of strings for equality, considering the order of elements.
func AreSlicesEqual(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	// Create copies to avoid mutating original slices
	copy1 := append([]string{}, slice1...)
	copy2 := append([]string{}, slice2...)

	// Sort both slices
	sort.Strings(copy1)
	sort.Strings(copy2)

	// Compare sorted slices
	for i := range copy1 {
		if copy1[i] != copy2[i] {
			return false
		}
	}
	return true
}
