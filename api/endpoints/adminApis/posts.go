package adminApis

import (
	"github.com/gin-gonic/gin"
)

func IncludePosts(router *gin.RouterGroup) {
	postsRouter := router.Group("/posts")
	{
		postsRouter.GET("/test", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"test": "Posts Test"}) })

	}

}
