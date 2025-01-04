package handlers

import (
	"net/http"
	"strings"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/api/service"
	"github.com/gin-gonic/gin"
)

type SubdomainHandler struct {
	SubdomainService *service.SubdomainService
}

func NewSubdomainHandler() *SubdomainHandler {
	return &SubdomainHandler{SubdomainService: service.NewSubdomainService()}
}

// GetSubdomainsByScopeHandler handles requests to retrieve subdomains by scope
func (h *SubdomainHandler) GetSubdomainsByScopeHandler(c *gin.Context) {
	scope := c.Param("domain") // Get 'scope' query parameter
	if scope == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Domain is required"})
		return
	}

	subdomains, err := h.SubdomainService.GetSubdomainsByScope(scope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subdomains"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subdomains": subdomains})
}
func (h *SubdomainHandler) GetSubdomainsByProgramName(c *gin.Context) {
	programName := c.Param("programname") // Get path parameter "programname"
	jsonOutput := c.Query("json")         // Get query parameter "json"

	if programName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "programname is required"})
		return
	}

	subdomains, err := h.SubdomainService.GetSubdomainsByProgramName(programName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve subdomains", "details": err.Error()})
		return
	}

	if jsonOutput == "true" {
		// Return JSON format
		c.JSON(http.StatusOK, gin.H{"subdomains": subdomains})
	} else {
		// Return plain text with newlines
		c.String(http.StatusOK, "%s", formatStringToAddBackSlashN(extractSubdomainNames(subdomains)))
	}
}

// Helper function to extract SubDomain names
func extractSubdomainNames(subdomains []model.Subdomain) []string {
	names := make([]string, len(subdomains))
	for i, sub := range subdomains {
		names[i] = sub.SubDomain
	}
	return names
}

// Helper function to format subdomains as newline-separated string
func formatStringToAddBackSlashN(obj []string) string {
	return strings.Join(obj, "\n")
}

func (h *SubdomainHandler) GetAllSubdomain(c *gin.Context) {
	jsonOutput := c.Query("json")
	subdomains, err := h.SubdomainService.GetAllSubdomain()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve subdomains", "details": err.Error()})
		return
	}

	if jsonOutput == "true" {
		// Return JSON format
		c.JSON(http.StatusOK, gin.H{"subdomains": subdomains})
	} else {
		// Return plain text with newlines
		c.String(http.StatusOK, "%s", formatStringToAddBackSlashN(extractSubdomainNames(subdomains)))
	}
}
