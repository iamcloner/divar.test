package adminApis

import (
	"divar.ir/api/repositories/adminRepositories"
	"divar.ir/api/repositories/authenticationRepositories"
	"divar.ir/api/repositories/userRepositories"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IncludeUsers(router *gin.RouterGroup) {
	usersRouter := router.Group("/users")
	{
		usersRouter.GET("/test", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"test": "Users test"}) })
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
		usersRouter.GET("/:userId/activeSessions", func(ctx *gin.Context) {
			userId := ctx.Param("userId")
			userIdObj, err := primitive.ObjectIDFromHex(userId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": "Invalid User ID"})
				return
			}
			res, err := adminRepositories.GetUserActiveSession(userIdObj)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, res)

		})
		usersRouter.GET("/:userId/inactiveSessions", func(ctx *gin.Context) {
			userId := ctx.Param("userId")
			userIdObj, err := primitive.ObjectIDFromHex(userId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": "Invalid User ID"})
				return
			}
			res, err := adminRepositories.GetUserInactiveSession(userIdObj)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, res)

		})
		usersRouter.PATCH("/:userId/changePassword", func(ctx *gin.Context) {
			userId := ctx.Param("userId")
			userIdObj, err := primitive.ObjectIDFromHex(userId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": "Invalid User ID"})
				return
			}
			newPassword, exist := ctx.GetPostForm("newPassword")
			if !exist {
				ctx.JSON(400, "wrong Input")
				return
			}
			err = userRepositories.UpdatePassword(userIdObj, newPassword)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, "password updated")

		})
		usersRouter.PATCH("/:userId/promote", func(ctx *gin.Context) {
			userId := ctx.Param("userId")
			userIdObj, err := primitive.ObjectIDFromHex(userId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": "Invalid User ID"})
				return
			}
			err = adminRepositories.PromoteToAdmin(userIdObj)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, "user promoted to admin")
		})
		usersRouter.PATCH("/:userId/verify", func(ctx *gin.Context) {
			userId := ctx.Param("userId")
			userIdObj, err := primitive.ObjectIDFromHex(userId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": "Invalid User ID"})
				return
			}
			err = authenticationRepositories.VerifyUser(userIdObj)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, "password updated")
		})
		usersRouter.PATCH("/:userId/lock", func(ctx *gin.Context) {
			userId := ctx.Param("userId")
			userIdObj, err := primitive.ObjectIDFromHex(userId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": "Invalid User ID"})
				return
			}
			err = adminRepositories.LockUser(userIdObj)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, "password updated")
		})
		usersRouter.PATCH("/:userId/unlock", func(ctx *gin.Context) {
			userId := ctx.Param("userId")
			userIdObj, err := primitive.ObjectIDFromHex(userId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": "Invalid User ID"})
				return
			}
			err = adminRepositories.UnlockUser(userIdObj)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, "password updated")
		})
		usersRouter.PATCH("/:userId/ban", func(ctx *gin.Context) {
			userId := ctx.Param("userId")
			userIdObj, err := primitive.ObjectIDFromHex(userId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": "Invalid User ID"})
				return
			}
			err = adminRepositories.BanUser(userIdObj)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, "password updated")
		})
		usersRouter.PATCH("/:userId/unban", func(ctx *gin.Context) {
			userId := ctx.Param("userId")
			userIdObj, err := primitive.ObjectIDFromHex(userId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": "Invalid User ID"})
				return
			}
			err = adminRepositories.UnbanUser(userIdObj)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, "password updated")
		})
		usersRouter.PATCH("/:userId/lock", func(ctx *gin.Context) {
			userId := ctx.Param("userId")
			userIdObj, err := primitive.ObjectIDFromHex(userId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": "Invalid User ID"})
				return
			}
			err = adminRepositories.LockUser(userIdObj)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, "password updated")
		})
		usersRouter.PATCH("/:userId/session/:sessionId/terminate", func(ctx *gin.Context) {
			userId := ctx.Param("userId")
			userIdObj, err := primitive.ObjectIDFromHex(userId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": "Invalid User ID"})
				return
			}
			sessionId := ctx.Param("sessionId")
			sessionIdObj, err := primitive.ObjectIDFromHex(sessionId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": "Invalid session ID"})
				return
			}
			err = adminRepositories.TerminateUserSession(userIdObj, sessionIdObj)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, "password updated")
		})
		usersRouter.PATCH("/:userId/session/terminateAll", func(ctx *gin.Context) {
			userId := ctx.Param("userId")
			userIdObj, err := primitive.ObjectIDFromHex(userId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": "Invalid User ID"})
				return
			}
			err = adminRepositories.TerminateAllSessionUser(userIdObj)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, "password updated")
		})
	}

}
