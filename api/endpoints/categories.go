package endpoints

import (
	"divar.ir/api/repositories"
	"github.com/gin-gonic/gin"
	"strconv"
)

func IncludeCategories(router *gin.Engine) {
	categoriesRouter := router.Group("/categories")
	{
		categoriesRouter.GET("/test", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"test": "categories"}) })
		categoriesRouter.GET("/", func(ctx *gin.Context) {
			result, err := repositories.GetCategories(0)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			}
			ctx.JSON(200, result)
		})
		categoriesRouter.GET("/:code", func(ctx *gin.Context) {
			code, err := strconv.Atoi(ctx.Param("code"))
			if err != nil {
				ctx.JSON(500, "invalid code")
				return
			}
			result, err := repositories.GetCategories(code)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			}
			ctx.JSON(200, result)
		})

	}

}
