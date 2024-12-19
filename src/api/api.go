package api

import (
	"github.com/bahramiofficial/watchtower/src/api/routers"
	"github.com/gin-gonic/gin"
)

func InitServer() {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	v1 := r.Group("api/v1")
	{
		testGroup := v1.Group("/test")
		routers.TestRouter(testGroup)
	}

	r.Run(":8000")
}
