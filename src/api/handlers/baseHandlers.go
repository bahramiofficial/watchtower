package handlers

import (
	"strconv"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/api/service"
	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
	BaseService *service.BaseService
}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{BaseService: service.NewBaseService()}
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

// Helper function to format subdomain list as a plain text response
func formatLiveSubdomainList(subdomains []model.LiveSubdomain) string {
	result := ""
	for _, sub := range subdomains {
		result += sub.Subdomain + "\n"
	}
	return result
}

func formatHttpSubdomainList(subdomains []model.Http) string {
	result := ""
	for _, sub := range subdomains {
		result += sub.SubDomain + "\n"
	}
	return result
}
