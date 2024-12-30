package endpoints

import (
	"divar.ir/api/repositories/postsRepositories"
	"github.com/gin-gonic/gin"
)

func IncludePosts(router *gin.Engine) {
	postsRouter := router.Group("/posts")
	{
		postsRouter.GET("/test", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"test": "posts"}) })
		postsRouter.GET("/", func(ctx *gin.Context) {
			areaCode := ctx.Query("areaCode")
			categoryCode := ctx.Query("categoryCode")
			result, err := postsRepositories.GetPosts(ctx, areaCode, categoryCode)
			if err != nil {
				return
			}
			ctx.JSON(200, result)
		})

	}

}
