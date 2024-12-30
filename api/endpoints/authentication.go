package endpoints

import (
	"divar.ir/api/repositories/authenticationRepositories"
	"github.com/gin-gonic/gin"
)

func IncludeAuthentication(router *gin.Engine) {
	loginRouter := router.Group("/auth")
	{
		loginRouter.GET("/test", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"test": "posts"}) })
		loginRouter.POST("/login", func(ctx *gin.Context) {
			result, err := authenticationRepositories.Login(ctx)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			} else {
				ctx.JSON(200, result)
			}
		})
		loginRouter.POST("/verify", func(ctx *gin.Context) {
			result, err := authenticationRepositories.Verify(ctx)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			} else {
				ctx.JSON(200, result)
			}
		})
		loginRouter.POST("/refresh", func(ctx *gin.Context) {
			result, err := authenticationRepositories.Refresh(ctx)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			} else {
				ctx.JSON(200, result)
			}
		})
		loginRouter.POST("/logout", func(ctx *gin.Context) {
			err := authenticationRepositories.Logout(ctx)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			} else {
				ctx.JSON(200, "Logout successfully")
			}
		})

	}
	registerRouter := router.Group("/register")
	{
		registerRouter.GET("/test", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"test": "posts"}) })
		registerRouter.POST("/", func(ctx *gin.Context) {
			result, err := authenticationRepositories.Register(ctx)
			if err != nil {
				return
			} else {
				ctx.JSON(200, result)
			}
		})

	}

}
