package endpoints

import (
	"divar.ir/api/repositories"
	"github.com/gin-gonic/gin"
)

func IncludeCategories(router *gin.Engine) {
	categoriesRouter := router.Group("/categories")
	{
		categoriesRouter.GET("/test", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"test": "categories"}) })
		categoriesRouter.GET("/", func(ctx *gin.Context) {
			result, err := repositories.GetCategories(ctx, "")
			if err != nil {
				return
			}
			ctx.JSON(200, result)
		})
		categoriesRouter.GET("/:code", func(ctx *gin.Context) {
			code := ctx.Param("code")
			result, err := repositories.GetCategories(ctx, code)
			if err != nil {
				return
			}
			ctx.JSON(200, result)
		})

	}

}
