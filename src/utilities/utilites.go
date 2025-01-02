package utilities

import (
	"strings"
	"time"
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
