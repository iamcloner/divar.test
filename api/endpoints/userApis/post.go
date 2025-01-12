package userApis

import (
	"divar.ir/api/repositories/userRepositories"
	"divar.ir/schema"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
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
			postInfo.CreatedAt = time.Now()
			postInfo.UpdatedAt = time.Now()
			postInfo.Status = "Pending"
			postInfo.LastAction = "Added by user"
			err = postInfo.Validate()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			}
			postInfo.UserId = userIdObj
			err = userRepositories.AddPost(postInfo)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			}
			ctx.JSON(200, "posted successfully")
		})
		postsRouter.PUT("/:post_id", func(ctx *gin.Context) {
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
			postID := ctx.Param("post_id")
			if postID == "" {
				ctx.JSON(http.StatusBadRequest, "post id is required")
				return
			}
			postInfo.ID, err = primitive.ObjectIDFromHex(postID)
			if err != nil {
				ctx.JSON(400, "Invalid post id")
				return
			}
			err = ctx.ShouldBind(&postInfo)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, "Invalid inputs")
				return
			}
			postInfo.UserId = userIdObj
			postInfo.UpdatedAt = time.Now()
			postInfo.Status = "Pending"
			postInfo.LastAction = "Updated by user"
			err = postInfo.Validate()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			}

			err = userRepositories.UpdatePost(postInfo)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			}
			ctx.JSON(200, "updated successfully")
		})
		postsRouter.DELETE("/:post_id", func(ctx *gin.Context) {
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
			postID := ctx.Param("post_id")
			if postID == "" {
				ctx.JSON(http.StatusBadRequest, "post id is required")
				return
			}
			postIdObj, err := primitive.ObjectIDFromHex(postID)
			if err != nil {
				ctx.JSON(400, "Invalid post id")
				return
			}

			err = userRepositories.DeletePost(userIdObj, postIdObj)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			}
			ctx.JSON(200, "Deleted successfully")
		})
	}

}
