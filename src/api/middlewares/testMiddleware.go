package middlewares

import (
	"github.com/didip/tollbooth"
 
	"github.com/gin-gonic/gin"
)

func LimitByRequest() gin.HandlerFunc{
	 // Create a request limiter per handler.
	 lmt := tollbooth.NewLimiter(1, nil)

	 return func(c *gin.Context){
		err := tollbooth.LimitByRequest(lmt, c.Writer,c.Request)
		if err!= nil{
            c.AbortWithStatusJSON(429, gin.H{"error": "Too Many Requests"})
            return
        }else{
			c.Next()
		}
	 }	 
	 
}