package main

import (
	"github.com/gin-gonic/gin"
	"github.com/happylusn/lithot-gin/example/configuration"
	"github.com/happylusn/lithot-gin/example/controllers"
	"github.com/happylusn/lithot-gin/lithot"
	"net/http"
)

func MyErrorHandle(c *gin.Context, err interface{}) {
	if err != nil {
		if s, ok := err.(string); ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": s})
		} else {
			if e, ok := err.(error); ok {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": e})
			}
		}
	}
}

func main() {
	lithot.NewLithot().
		Configure(configuration.NewTestConfig(), configuration.NewGormConfig()).
		Mount("/v1", controllers.NewUserController()).
		SetErrorHandle(MyErrorHandle).
		Launch()
}
