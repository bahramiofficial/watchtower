package routers

import (
	"github.com/bahramiofficial/watchtower/src/api/handlers"
	"github.com/gin-gonic/gin"
)

func ProgramRouter(r *gin.RouterGroup) {
	programHandler := handlers.NewProgramHandler()

	r.GET("/", programHandler.GetAllProgramsHandler)

}
