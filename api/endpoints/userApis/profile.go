package userApis

import (
	"divar.ir/api/repositories/userRepositories"
	"divar.ir/internal/email"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

func IncludeProfile(router *gin.RouterGroup) {
	profileRouter := router.Group("/profile")
	{
		profileRouter.GET("/test", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"test": "Profile"}) })
		profileRouter.GET("/", func(ctx *gin.Context) {
			userId, exist := ctx.Get("userId")
			if !exist {
				ctx.JSON(401, "Authentication Required")
				return
			}
			userIdStr := userId.(string)
			profile, err := userRepositories.GetProfile(userIdStr)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			}
			ctx.JSON(200, profile)
		})
		profileRouter.PATCH("/updateName", func(ctx *gin.Context) {

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
			newName, exist := ctx.GetPostForm("newName")
			if !exist {
				ctx.JSON(400, "wrong Input")
				return
			}
			if len(newName) < 5 || len(newName) > 16 {
				ctx.JSON(400, "Name must between 5 and 16 characters")
				return
			}
			err = userRepositories.UpdateProfileName(userIdObj, newName)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			}
			ctx.JSON(200, "name updated successfully")
		})
		profileRouter.PUT("/updateEmail", func(ctx *gin.Context) {

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
			oldEmail, exist := ctx.GetPostForm("oldEmail")
			if !exist || strings.Count(oldEmail, "@") > 1 {
				ctx.JSON(400, "wrong Input")
				return
			}

			err = userRepositories.UpdateProfileEmail(userIdObj, oldEmail)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			}
			ctx.JSON(200, "enter code.check your Email.")
		})
		profileRouter.PATCH("/confirmUpdateEmail", func(ctx *gin.Context) {

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
			verifyCode, exist := ctx.GetPostForm("code")
			if !exist || len(verifyCode) != 5 {
				ctx.JSON(400, "wrong Input")
				return
			}
			newEmail, exist := ctx.GetPostForm("newEmail")
			if !exist || !email.IsVail(newEmail) {
				ctx.JSON(400, "wrong Input")
				return
			}

			err = userRepositories.ConfirmUpdateProfileEmail(userIdObj, verifyCode, newEmail)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			}
			ctx.JSON(200, "email updated successfully.")
		})
		profileRouter.POST("/updateAvatar", func(ctx *gin.Context) {

			ctx.JSON(200, "profile")
		})
		profileRouter.DELETE("/", func(ctx *gin.Context) {
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
			err = userRepositories.DeleteProfile(userIdObj)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			}
			ctx.JSON(200, "enter code.check your Email.")
		})
		profileRouter.PATCH("/confirmDelete", func(ctx *gin.Context) {
			userId, exist := ctx.Get("userId")
			if !exist {
				ctx.JSON(401, "Authentication Required")
				return
			}
			verifyCode, exist := ctx.GetPostForm("code")
			if !exist || len(verifyCode) != 5 {
				ctx.JSON(400, "wrong Input")
				return
			}
			userIdObj, err := primitive.ObjectIDFromHex(userId.(string))
			if err != nil {
				ctx.JSON(401, "Authentication Failed")
				return
			}
			err = userRepositories.ConfirmDeleteProfile(userIdObj, verifyCode)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			}
			ctx.JSON(200, "successfully deleted.")
		})
		profileRouter.PATCH("/updatePassword", func(ctx *gin.Context) {

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
			oldPassword, exist := ctx.GetPostForm("oldPassword")
			if !exist {
				ctx.JSON(400, "wrong Input")
				return
			}
			newPassword, exist := ctx.GetPostForm("newPassword")
			if !exist {
				ctx.JSON(400, "wrong Input")
				return
			}
			check, err := userRepositories.CheckOldPassword(userIdObj, oldPassword)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			}
			if !check {
				ctx.JSON(400, "wrong Password")
				return
			}
			err = userRepositories.UpdatePassword(userIdObj, newPassword)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			}
			ctx.JSON(200, "password updated successfully.")
		})
	}

}
