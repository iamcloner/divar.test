package userApis

import (
	"divar.ir/api/repositories/userRepositories"
	"github.com/gin-gonic/gin"
)

func IncludeProfile(router *gin.RouterGroup) {
	profileRouter := router.Group("/profile")
	{
		profileRouter.GET("/test", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"test": "Profile"}) })
		profileRouter.GET("/", func(ctx *gin.Context) {
			userId, exist := ctx.Get("userId")
			if !exist {
				ctx.JSON(401, "Authentication Required")
			}
			userIdStr := userId.(string)
			profile, err := userRepositories.GetProfile(userIdStr)
			if err != nil {
				ctx.JSON(500, err.Error())
			}
			ctx.JSON(200, profile)
		})

	}

}
