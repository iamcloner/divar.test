package endpoints

import (
	"divar.ir/api/repositories"
	"github.com/gin-gonic/gin"
)

func IncludePost(router *gin.Engine) {
	postsRouter := router.Group("/post")
	{
		postsRouter.GET("/test", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"test": "posts"}) })
		postsRouter.GET("/:id", func(ctx *gin.Context) {
			id := ctx.Param("id")
			result, err := repositories.GetPost(ctx, id)
			if err != nil {
				return
			}
			ctx.JSON(200, result)
		})

		postsRouter.GET("/:id/phone", func(ctx *gin.Context) {
			id := ctx.Param("id")
			result, err := repositories.GetPostPhone(ctx, id)
			if err != nil {
				return
			}
			ctx.JSON(200, result)
		})
	}

}
