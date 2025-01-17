package adminApis

import (
	"divar.ir/api/repositories/adminRepositories"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IncludePosts(router *gin.RouterGroup) {
	postsRouter := router.Group("/posts")
	{
		postsRouter.GET("/test", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"test": "Posts Test"}) })
		postsRouter.GET("/getPendingPosts", func(ctx *gin.Context) {

			posts, err := adminRepositories.GetPendingPosts()
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, posts)
		})
		postsRouter.POST("/:postId/Verify", func(ctx *gin.Context) {
			postId := ctx.Param("postId")
			postIdObj, err := primitive.ObjectIDFromHex(postId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": "Invalid post ID"})
				return
			}
			err = adminRepositories.VerifyPost(postIdObj)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, "post verified")
		})
		postsRouter.POST("/:postId/reject", func(ctx *gin.Context) {
			postId := ctx.Param("postId")
			postIdObj, err := primitive.ObjectIDFromHex(postId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": "Invalid post ID"})
				return
			}
			err = adminRepositories.RejectPost(postIdObj)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, "post verified")
		})
	}

}
