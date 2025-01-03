package api

import (
	"github.com/bahramiofficial/watchtower/src/api/middlewares"
	"github.com/bahramiofficial/watchtower/src/api/routers"
	"github.com/bahramiofficial/watchtower/src/api/validations"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitServer() {
	r := gin.New()
	validation, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		validation.RegisterValidation("mobile", validations.IranianMobileNumberValidator, true)
	}
	r.Use(gin.Logger(), gin.Recovery(), middlewares.LimitByRequest())

	v1 := r.Group("api/v1")
	{
		// testGroup := v1.Group("/test")
		// routers.TestRouter(testGroup)
		programGroup := v1.Group("/program")
		routers.ProgramRouter(programGroup)

		subdomianGroup := v1.Group("/subdomian")
		routers.SubDomainRouter(subdomianGroup)
	}

	r.Run(":8000")
}
