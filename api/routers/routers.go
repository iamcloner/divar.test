package routers

import (
	"divar.ir/api/endpoints"
	"github.com/gin-gonic/gin"
)

func IncludeRouters(router *gin.Engine) {
	endpoints.IncludeCities(router)
	endpoints.IncludeCategories(router)
	endpoints.IncludePosts(router)
	endpoints.IncludePost(router)
	endpoints.IncludeAuthentication(router)

	IncludeUserRouters(router)
}
