package adminApis

import (
	"divar.ir/internal/mongodb"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func IncludeUsers(router *gin.RouterGroup) {
	usersRouter := router.Group("/users")
	{
		usersRouter.GET("/test", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"test": "Cities"}) })
		usersRouter.GET("/", func(ctx *gin.Context) {

			handler, err := mongodb.GetMongoDBHandler()

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Operation Failed (0x0001)"})
				return
			}

			result, err := handler.FindMany("cities", bson.M{}, bson.M{"_id": 0})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Operation Failed (0x0002)"})
				return
			}

			ctx.JSON(200, result)
		})
	}

}
