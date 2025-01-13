package adminApis

import (
	"divar.ir/api/repositories/adminRepositories"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IncludeUsers(router *gin.RouterGroup) {
	usersRouter := router.Group("/users")
	{
		usersRouter.GET("/test", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"test": "Cities"}) })
		usersRouter.GET("/", func(ctx *gin.Context) {

			res, err := adminRepositories.GetUsers()
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, gin.H{"users": res})

		})
		usersRouter.GET("/:userId", func(ctx *gin.Context) {
			userId := ctx.Param("userId")
			userIdObj, err := primitive.ObjectIDFromHex(userId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": "Invalid User ID"})
				return
			}
			res, err := adminRepositories.GetUser(userIdObj)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, res)

		})
		usersRouter.GET("/:userId/loginInfo", func(ctx *gin.Context) {
			userId := ctx.Param("userId")
			userIdObj, err := primitive.ObjectIDFromHex(userId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": "Invalid User ID"})
				return
			}
			res, err := adminRepositories.GetUserLoginInfo(userIdObj)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, res)

		})
	}

}
