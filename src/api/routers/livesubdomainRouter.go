package routers

import (
	"github.com/bahramiofficial/watchtower/src/api/handlers"
	"github.com/gin-gonic/gin"
)

func LiveSubDomainRouter(r *gin.RouterGroup) {
	livesubdomainHandler := handlers.NewLiveSubdomainHandler()

	r.GET("/", livesubdomainHandler.LiveSubdomainsHandler)
	r.GET("/details/:subdomain ", livesubdomainHandler.GetLiveSubdomainDetailHandler)

}
