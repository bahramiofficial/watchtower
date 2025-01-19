package handlers

import (
	"net/http"

	"github.com/bahramiofficial/watchtower/src/api/service"
	"github.com/gin-gonic/gin"
)

type LiveSubdomainHandler struct {
	LiveSubdomainService *service.LiveSubdomainService
}

func NewLiveSubdomainHandler() *LiveSubdomainHandler {
	return &LiveSubdomainHandler{LiveSubdomainService: service.NewLiveSubdomainService()}
}

// LiveSubdomainsHandler handles the /api/livesubdomains endpoint
func (h *LiveSubdomainHandler) LiveSubdomainsHandler(c *gin.Context) {
	filter := service.LiveSubdomainFilter{
		ProgramName: c.Query("program"),
		Scope:       c.Query("scope"),
		Provider:    c.Query("provider"),
		Fresh:       c.DefaultQuery("fresh", "false") == "true",
		Count:       c.DefaultQuery("count", "false") == "true",
		Limit:       parseQueryInt(c, "limit", 1000),
		Page:        parseQueryInt(c, "page", 1),
	}
	jsonOutput := c.DefaultQuery("json", "false") == "true"

	livesubdomains, count, err := h.LiveSubdomainService.GetLiveSubdomains(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve subdomains", "details": err.Error()})
		return
	}

	if filter.Count {
		c.JSON(http.StatusOK, gin.H{"count": count})
		return
	}

	if len(livesubdomains) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No subdomains found"})
		return
	}

	if jsonOutput {
		c.JSON(http.StatusOK, livesubdomains)
	} else {
		c.String(http.StatusOK, formatLiveSubdomainList(livesubdomains))
	}
}

// GetLiveSubdomainDetailHandler handles the /api/livesubdomains/details/:subdomain endpoint
func (h *LiveSubdomainHandler) GetLiveSubdomainDetailHandler(c *gin.Context) {
	subdomain := c.Param("subdomain")
	subdomainObj, err := h.LiveSubdomainService.GetSingleLiveSubdomainBySubDomain(subdomain)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found", "subdomain": subdomain})
		return
	}

	c.JSON(http.StatusOK, subdomainObj)
}
