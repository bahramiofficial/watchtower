package routers

import (
	"github.com/bahramiofficial/watchtower/src/api/handlers"
	"github.com/gin-gonic/gin"
)

func HttpRouter(r *gin.RouterGroup) {
	httpHandler := handlers.NewHttpHandler()

	r.GET("/", httpHandler.HttpHandler)
	r.GET("/details/:subdomain ", httpHandler.GetHttpDetailHandler)

}
