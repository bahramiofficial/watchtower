package handlers

import (
	"net/http"

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
