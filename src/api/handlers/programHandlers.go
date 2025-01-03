package handlers

import (
	"net/http"

	"github.com/bahramiofficial/watchtower/src/api/service"
	"github.com/gin-gonic/gin"
)

type ProgramHandler struct {
	ProgramService *service.ProgramService
}

func NewProgramHandler() *ProgramHandler {
	return &ProgramHandler{ProgramService: service.NewProgramService()}
}

// GetAllProgramsHandler handles requests to retrieve all programs
func (h *ProgramHandler) GetAllProgramsHandler(c *gin.Context) {
	programs, err := h.ProgramService.GetAllPrograms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve programs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"programs": programs})
}
