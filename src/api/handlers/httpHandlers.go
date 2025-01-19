package handlers

import (
	"net/http"

	"github.com/bahramiofficial/watchtower/src/api/service"
	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	HttpService *service.HttpService
}

func NewHttpHandler() *HttpHandler {
	return &HttpHandler{HttpService: service.NewHttpService()}
}

// SubdomainsHandler handles the /api/livesubdomains endpoint
func (h *HttpHandler) HttpHandler(c *gin.Context) {

	filter := service.HttpFilter{
		ProgramName: c.Query("program"),
		Scope:       c.Query("scope"),
		Provider:    c.Query("provider"),
		Title:       c.Query("title"),
		Status:      c.Query("status"),
		Tech:        c.Query("tech"),
		Latest:      c.DefaultQuery("latest", "false") == "true",
		Fresh:       c.DefaultQuery("fresh", "false") == "true",
		Count:       c.DefaultQuery("count", "false") == "true",
		Limit:       parseQueryInt(c, "limit", 1000),
		Page:        parseQueryInt(c, "page", 1),
		// TODO: Add more filters as needed (e.g., created_at, updated_at, etc.)
	}

	jsonOutput := c.DefaultQuery("json", "false") == "true"

	httpobj, count, err := h.HttpService.GetHttp(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve subdomains", "details": err.Error()})
		return
	}

	if filter.Count {
		c.JSON(http.StatusOK, gin.H{"count": count})
		return
	}

	if len(httpobj) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No subdomains found"})
		return
	}

	if jsonOutput {
		c.JSON(http.StatusOK, httpobj)
	} else {
		c.String(http.StatusOK, formatHttpSubdomainList(httpobj))
	}
}

// GetSubdomainDetailHandler handles the /api/livesubdomains/details/:subdomain endpoint
func (h *HttpHandler) GetHttpDetailHandler(c *gin.Context) {
	subdomain := c.Param("subdomain")
	httpObj, err := h.HttpService.GetSingleHttpBySubDomain(subdomain)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found", "subdomain": subdomain})
		return
	}

	c.JSON(http.StatusOK, httpObj)
}
