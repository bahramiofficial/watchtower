package routers

import (
	"github.com/bahramiofficial/watchtower/src/api/handlers"
	"github.com/gin-gonic/gin"
)

func ProgramRouter(r *gin.RouterGroup) {
	programHandler := handlers.NewProgramHandler()

	r.GET("/all", programHandler.GetAllProgramsHandler)
	r.GET("/:programname", programHandler.GetSingleProgramHandler)
	r.GET("/:programname/delete", programHandler.DeleteProgramHandler)

}
