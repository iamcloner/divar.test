package postsRepositories

import (
	"divar.ir/internal/mongodb"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func GetPost(ctx *gin.Context, id string) (bson.M, error) {
	handler, err := mongodb.GetMongoDBHandler()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Operation Failed (0x0001)"})
		return nil, err
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return nil, err
	}

	var result bson.M
	err = handler.FindOne("posts", bson.M{"_id": objID},
		bson.M{"phone": 0, "userId": 0}).Decode(&result)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Operation Failed (0x0002)"})
		return nil, err
	}
	return result, nil

}

func GetPostPhone(ctx *gin.Context, id string) (bson.M, error) {
	handler, err := mongodb.GetMongoDBHandler()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Operation Failed (0x0001)"})
		return nil, err
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return nil, err
	}

	var result bson.M
	err = handler.FindOne("posts", bson.M{"_id": objID}, bson.M{"userId": 1}).Decode(&result)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Operation Failed (0x0002)"})
		return nil, err
	}
	err = handler.FindOne("users", bson.M{"_id": result["userId"]}, bson.M{"phone": 1}).Decode(&result)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Operation Failed (0x0002)"})
		return nil, err
	}
	return bson.M{"phone": result["phone"], "_id": id}, nil

}
