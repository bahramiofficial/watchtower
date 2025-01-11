package handlers

import (
	"net/http"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/api/service"
	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	HttpService *service.HttpService
}

func NewHttpHandler() *HttpHandler {
	return &HttpHandler{HttpService: service.NewHttpService()}
}

// GetHttpByScopeHandler handles requests to retrieve Http data by scope
func (h *HttpHandler) GetHttpByScopeHandler(c *gin.Context) {
	scope := c.Param("domain")
	if scope == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Domain is required"})
		return
	}

	httpData, err := h.HttpService.GetHttpByScope(scope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve HTTP data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"httpData": httpData})
}

// GetHttpByProgramNameHandler handles requests to retrieve Http data by program name
func (h *HttpHandler) GetHttpByProgramNameHandler(c *gin.Context) {
	programName := c.Param("programname")
	jsonOutput := c.Query("json")

	if programName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "programname is required"})
		return
	}

	httpData, err := h.HttpService.GetHttpByProgramName(programName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve HTTP data", "details": err.Error()})
		return
	}

	if jsonOutput == "true" {
		c.JSON(http.StatusOK, gin.H{"httpData": httpData})
	} else {
		c.String(http.StatusOK, "%s", formatStringToAddBackSlashN(extractHttpNames(httpData)))
	}
}

// Helper function to extract Http names
func extractHttpNames(httpData []model.Http) []string {
	names := make([]string, len(httpData))
	for i, http := range httpData {
		names[i] = http.SubDomain
	}
	return names
}
