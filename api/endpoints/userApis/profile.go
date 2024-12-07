package userApis

import (
	"github.com/gin-gonic/gin"
)

func IncludeProfile(router *gin.RouterGroup) {
	profileRouter := router.Group("/profile")
	{
		profileRouter.GET("/test", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"test": "posts"}) })
		profileRouter.GET("/", func(ctx *gin.Context) {

			ctx.JSON(200, "")
		})

	}

}
