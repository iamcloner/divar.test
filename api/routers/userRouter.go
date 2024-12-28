package routers

import (
	"divar.ir/api/endpoints/userApis"
	"divar.ir/api/middleware"
	"github.com/gin-gonic/gin"
)

func IncludeUserRouters(router *gin.Engine) {
	userRouter := router.Group("/user")
	{
		userRouter.Use(middleware.UserAuthentication())
		userApis.IncludeProfile(userRouter)
		userApis.IncludePost(userRouter)
	}

}
