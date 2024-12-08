package repositories

import (
	"divar.ir/internal/mongodb"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
)

func GetCategories(ctx *gin.Context, code string) (interface{}, error) {
	handler, err := mongodb.GetMongoDBHandler()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Operation Failed (0x0001)"})
		return nil, err
	}

	if code == "" {
		result, err := handler.FindMany("categories", bson.M{}, bson.M{"_id": 0, "subs.filters": 0})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Operation Failed (0x0002)"})
			return nil, err
		}
		return result, nil
	} else {
		codei, err := strconv.Atoi(code)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Id (0x0002)"})
			return nil, err
		}
		sub := codei % 100
		cat := codei / 100
		if cat < 1 {
			var tmp bson.M
			err := handler.FindOne("categories", bson.M{"code": codei}, bson.M{"_id": 0, "subs.filters": 0}).Decode(&tmp)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Operation Failed (0x0002)"})
				return nil, err
			}
			return tmp, nil

		} else {

			var tmp bson.M
			err := handler.FindOne("categories", bson.M{"code": cat, "subs.code": sub}, bson.M{"subs": bson.M{"$elemMatch": bson.M{"code": sub}}, "_id": 0}).Decode(&tmp)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Operation Failed (0x0002)"})
				return nil, err
			}

			tmp = (tmp["subs"].(bson.A)[0]).(bson.M)
			tmp["code"] = cat*100 + sub
			return tmp, nil
		}
	}

}
