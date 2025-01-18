package handlers

import (
	"net/http"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/api/service"
	"github.com/gin-gonic/gin"
)

type LiveSubdomainHandler struct {
	LiveSubdomainService *service.LiveSubdomainService
	SubdomainService     *service.SubdomainService
}

func NewLiveSubdomainHandler() *LiveSubdomainHandler {
	return &LiveSubdomainHandler{LiveSubdomainService: service.NewLiveSubdomainService(), SubdomainService: service.NewSubdomainService()}
}

type CombinedSubdomain struct {
	model.LiveSubdomain
	Providers model.StringArray `json:"providers"`
}

// GetLiveSubdomainsByScopeHandler handles requests to retrieve live subdomains by scope
func (h *LiveSubdomainHandler) GetSingleLiveSubdomainHandler(c *gin.Context) {
	subdomainInput := c.Param("subdomianlive")
	if subdomainInput == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Domain is required"})
		return
	}

	liveSubdomain, err := h.LiveSubdomainService.GetSingleLiveSubdomainBySubdomain(subdomainInput)
	if err != nil {
		// Handle error
	}

	subdomain, err := h.SubdomainService.GetSingleSubdomainBySubDomain(subdomainInput)
	if err != nil {
		// Handle error
	}

	// Assign Providers from Subdomain to LiveSubdomain
	combined := CombinedSubdomain{
		LiveSubdomain: liveSubdomain,
		Providers:     subdomain.Providers,
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve live subdomains"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"liveSubdomains": combined})
}

// todo edit code and ad service for get live domomain with provider
// GetLiveSubdomainsByScopeHandler handles requests to retrieve live subdomains by scope
func (h *LiveSubdomainHandler) GetLiveSubdomainWithProviderHandler(c *gin.Context) {
	provider := c.Param("provider")

	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Domain is required"})
		return
	}
	//todo edit code and ad service for get live domomain with provider
	// time =  time.Now(). - timedelta(hurse=12).
	// get all sub domain with provider  from subdmoain  and for on all
	// print sub domain
	// get livesubdomain  where subdomain = subsodmina and updatedat  bozrgtar az time
	// if live have value  string res  subdomain+\n
	//return res

	c.JSON(http.StatusOK, gin.H{"liveSubdomains": "liveSubdomains"})
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
func extractLiveSubdomainNames(liveSubdomains []model.LiveSubdomain) []string {
	names := make([]string, len(liveSubdomains))
	for i, liveSubdomain := range liveSubdomains {
		names[i] = liveSubdomain.Subdomain
	}
	return names
}
