package userApis

import (
	"divar.ir/api/repositories/userRepositories"
	"divar.ir/schema"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func IncludePost(router *gin.RouterGroup) {
	postsRouter := router.Group("/posts")
	{
		postsRouter.GET("/", func(ctx *gin.Context) {
			userId, exist := ctx.Get("userId")
			if !exist {
				ctx.JSON(401, "Authentication Required")
				return
			}
			userIdObj, err := primitive.ObjectIDFromHex(userId.(string))
			if err != nil {
				ctx.JSON(401, "Authentication Failed")
				return
			}
			profile, err := userRepositories.GetMyPosts(userIdObj)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			}
			ctx.JSON(200, profile)
		})
		postsRouter.POST("/", func(ctx *gin.Context) {
			userId, exist := ctx.Get("userId")
			if !exist {
				ctx.JSON(401, "Authentication Required")
				return
			}
			userIdObj, err := primitive.ObjectIDFromHex(userId.(string))
			if err != nil {
				ctx.JSON(401, "Authentication Failed")
				return
			}
			var postInfo schema.Posts
			err = ctx.ShouldBind(&postInfo)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, "Invalid inputs")
				return
			}
			err = userRepositories.CheckPostRequirements(postInfo)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			}

			err = userRepositories.AddPost(userIdObj, postInfo)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			}
			ctx.JSON(200, "posted successfully")
		})
	}

}
