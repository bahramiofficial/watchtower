package handlers

import (
	"net/http"
	"strings"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/api/service"
	"github.com/gin-gonic/gin"
)

type ProgramHandler struct {
	ProgramService *service.ProgramService
}

func NewProgramHandler() *ProgramHandler {
	return &ProgramHandler{ProgramService: service.NewProgramService()}
}

// Helper function to extract SubDomain names
func extractProgramNames(programs []model.Program) []string {
	names := make([]string, len(programs))
	for i, pro := range programs {
		names[i] = pro.ProgramName
	}
	return names
}

// Helper function to format subdomains as newline-separated string
func formatStringToAddBackSlashN(obj []string) string {
	return strings.Join(obj, "\n")
}

// GetAllProgramsHandler handles requests to retrieve all programs
func (h *ProgramHandler) GetAllProgramsHandler(c *gin.Context) {
	jsonOutput := c.Query("json")
	programs, err := h.ProgramService.GetAllPrograms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve programs"})
		return
	}

	if jsonOutput == "true" {
		// Return JSON format
		c.JSON(http.StatusOK, gin.H{"programs": programs})
	} else {
		// Return plain text with newlines
		c.String(http.StatusOK, "%s", formatStringToAddBackSlashN(extractProgramNames(programs)))
	}
}

func (h *ProgramHandler) GetSingleProgramsHandler(c *gin.Context) {
	jsonOutput := c.Query("json")
	programName := c.Param("programname")
	if programName == "" {
		if jsonOutput == "true" {
			// Return JSON format
			c.JSON(http.StatusBadRequest, gin.H{"error": "programname is required"})
		} else {
			// Return plain text with newlines
			c.String(http.StatusBadRequest, "programname is required")
		}
	}
	program, err := h.ProgramService.GetSinglePrograms(programName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve programs"})
		return
	}

	if jsonOutput == "true" {
		// Return JSON format
		c.JSON(http.StatusOK, gin.H{"programs": program})
	} else {
		// Return plain text with newlines
		c.String(http.StatusOK, "%s", program)
	}
}
