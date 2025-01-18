package routers

import (
	"github.com/bahramiofficial/watchtower/src/api/handlers"
	"github.com/gin-gonic/gin"
)

func SubDomainRouter(r *gin.RouterGroup) {
	subdomainHandler := handlers.NewSubdomainHandler()
	r.GET("/", subdomainHandler.SubdomainsHandler)
	r.GET("/details/:subdomain ", subdomainHandler.GetSubdomainDetailHandler)

}
