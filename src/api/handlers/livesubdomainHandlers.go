package handlers

import (
	"net/http"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/api/service"
	"github.com/gin-gonic/gin"
)

type LiveSubdomainHandler struct {
	LiveSubdomainService *service.LiveSubdomainService
}

func NewLiveSubdomainHandler() *LiveSubdomainHandler {
	return &LiveSubdomainHandler{LiveSubdomainService: service.NewLiveSubdomainService()}
}

// GetLiveSubdomainsByScopeHandler handles requests to retrieve live subdomains by scope
func (h *LiveSubdomainHandler) GetLiveSubdomainsByScopeHandler(c *gin.Context) {
	scope := c.Param("domain")
	if scope == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Domain is required"})
		return
	}

	liveSubdomains, err := h.LiveSubdomainService.GetLiveSubdomainsByScope(scope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve live subdomains"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"liveSubdomains": liveSubdomains})
}

// GetLiveSubdomainsByProgramNameHandler handles requests to retrieve live subdomains by program name
func (h *LiveSubdomainHandler) GetLiveSubdomainsByProgramNameHandler(c *gin.Context) {
	programName := c.Param("programname")
	jsonOutput := c.Query("json")

	if programName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "programname is required"})
		return
	}

	liveSubdomains, err := h.LiveSubdomainService.GetLiveSubdomainsByProgramName(programName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve live subdomains", "details": err.Error()})
		return
	}

	if jsonOutput == "true" {
		c.JSON(http.StatusOK, gin.H{"liveSubdomains": liveSubdomains})
	} else {
		c.String(http.StatusOK, "%s", formatStringToAddBackSlashN(extractLiveSubdomainNames(liveSubdomains)))
	}
}

// Helper function to extract LiveSubdomain names
func extractLiveSubdomainNames(liveSubdomains []model.LiveSubdomains) []string {
	names := make([]string, len(liveSubdomains))
	for i, liveSubdomain := range liveSubdomains {
		names[i] = liveSubdomain.SubDomain
	}
	return names
}
