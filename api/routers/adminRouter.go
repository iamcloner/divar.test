package routers

import (
	"divar.ir/api/endpoints/adminApis"
	"divar.ir/api/middleware"
	"github.com/gin-gonic/gin"
)

func IncludeAdminRouters(router *gin.Engine) {
	adminRouter := router.Group("/admin")
	{
		adminRouter.Use(middleware.AdminAuthentication())
		adminApis.IncludeUsers(adminRouter)
		adminApis.IncludePosts(adminRouter)
	}

}
