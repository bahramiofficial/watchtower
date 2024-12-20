package routers

import (
	"github.com/bahramiofficial/watchtower/src/api/handlers"
	"github.com/gin-gonic/gin"
)

func TestRouter(r *gin.RouterGroup) {
	testHandler := handlers.NewTestHandler()

	r.GET("/", testHandler.Test)
	r.GET(":id", testHandler.TestById)
	r.GET("/a/*optinals", testHandler.TestById)

}
