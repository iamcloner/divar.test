package repositories

import (
	"divar.ir/internal/mongodb"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
)

func GetPosts(ctx *gin.Context, areaCode string, categoryCode string) ([]bson.M, error) {
	handler, err := mongodb.NewMongoDBHandler()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Operation Failed (0x0001)"})
		return nil, err
	}

	codei, err := strconv.Atoi(areaCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Area Code (0x0002)"})
		return nil, err
	}
	zone := codei % 100
	area := codei % 10000 / 100
	city := codei / 10000

	var areaCodeLow, areaCodeHigh int
	var cateCodeLow, cateCodeHigh int

	if city < 1 && area < 1 {
		areaCodeHigh = zone * 10000
		areaCodeLow = zone*10000 + 9999
	} else if city < 1 && area > 1 {
		areaCodeHigh = area*10000 + zone*100
		areaCodeLow = area*10000 + zone*100 + 99

	} else {
		areaCodeLow = codei
		areaCodeHigh = codei
	}

	if categoryCode == "" {
		var result []bson.M
		result, err := handler.FindMany("posts", bson.M{"area-code": bson.M{"$gte": areaCodeHigh, "$lte": areaCodeLow}},
			bson.M{"title": 1, "images": 1, "area-code": 1, "price": 1, "used": 1})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Operation Failed (0x0002)"})
			return nil, err
		}
		return result, nil
	} else {
		codei, err = strconv.Atoi(categoryCode)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Area Code (0x0002)"})
			return nil, err
		}
		if codei/100 > 1 {
			cateCodeLow = codei
			cateCodeHigh = codei
		} else {
			cateCodeHigh = codei * 100
			cateCodeLow = codei*100 + 99
		}

		var result []bson.M
		result, err = handler.FindMany("posts", bson.M{
			"area-code":     bson.M{"$gte": areaCodeHigh, "$lte": areaCodeLow},
			"category-code": bson.M{"$gte": cateCodeHigh, "$lte": cateCodeLow}},
			bson.M{"title": 1, "images": 1, "area-code": 1, "price": 1, "used": 1})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Operation Failed (0x0002)"})
			return nil, err
		}
		return result, nil
	}

}
