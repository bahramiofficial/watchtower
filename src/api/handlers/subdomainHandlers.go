package handlers

import (
	"net/http"
	"strconv"

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

// SubdomainsHandler handles the /api/subdomains endpoint
func (h *SubdomainHandler) SubdomainsHandler(c *gin.Context) {
	filter := service.SubdomainFilter{
		ProgramName: c.Query("program"),
		Scope:       c.Query("scope"),
		Provider:    c.Query("provider"),
		Fresh:       c.DefaultQuery("fresh", "false") == "true",
		Count:       c.DefaultQuery("count", "false") == "true",
		Limit:       parseQueryInt(c, "limit", 1000),
		Page:        parseQueryInt(c, "page", 1),
	}
	jsonOutput := c.DefaultQuery("json", "false") == "true"

	subdomains, count, err := h.SubdomainService.GetSubdomains(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve subdomains", "details": err.Error()})
		return
	}

	if filter.Count {
		c.JSON(http.StatusOK, gin.H{"count": count})
		return
	}

	if len(subdomains) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No subdomains found"})
		return
	}

	if jsonOutput {
		c.JSON(http.StatusOK, subdomains)
	} else {
		c.String(http.StatusOK, formatSubdomainList(subdomains))
	}
}

// GetSubdomainDetailHandler handles the /api/subdomains/details/:subdomain endpoint
func (h *SubdomainHandler) GetSubdomainDetailHandler(c *gin.Context) {
	subdomain := c.Param("subdomain")
	subdomainObj, err := h.SubdomainService.GetSingleSubdomainBySubDomain(subdomain)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found", "subdomain": subdomain})
		return
	}

	c.JSON(http.StatusOK, subdomainObj)
}

// Helper function to parse query parameters as integers with a default value
func parseQueryInt(c *gin.Context, key string, defaultValue int) int {
	if val, ok := c.GetQuery(key); ok {
		if intValue, err := strconv.Atoi(val); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Helper function to format subdomain list as a plain text response
func formatSubdomainList(subdomains []model.Subdomain) string {
	result := ""
	for _, sub := range subdomains {
		result += sub.SubDomain + "\n"
	}
	return result
}
