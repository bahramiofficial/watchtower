package routers

import (
	"github.com/bahramiofficial/watchtower/src/api/handlers"
	"github.com/gin-gonic/gin"
)

func SubDomainRouter(r *gin.RouterGroup) {
	subdomainHandler := handlers.NewSubdomainHandler()

	r.GET("/all", subdomainHandler.GetAllSubdomain)
	r.GET("/domain/:domain", subdomainHandler.GetSubdomainsByScopeHandler)
	r.GET("/program/:programname", subdomainHandler.GetSubdomainsByProgramName)

}
