package postsRepositories

import (
	"divar.ir/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

func AreaCodeExist(areaCode int) bool {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return false
	}
	if areaCode < 100 {
		var result bson.M
		err = handler.FindOne("cities", bson.M{"code": areaCode}, bson.M{}).Decode(&result)
		if err != nil {
			return false
		}
	} else if areaCode > 1000 && areaCode < 9999 {
		city := areaCode / 100
		area := areaCode % 100
		var result bson.M
		err = handler.FindOne("cities", bson.M{"code": city, "area.$.code": area}, bson.M{}).Decode(&result)
		if err != nil {
			return false
		}
	} else if areaCode > 100000 && areaCode < 999999 {
		zone := areaCode % 100
		area := areaCode % 10000 / 100
		city := areaCode / 10000

		var result bson.M
		err = handler.FindOne("cities", bson.M{"code": city, "area.code": area, "area.zones.code": zone}, bson.M{}).Decode(&result)
		if err != nil {
			return false
		}
	} else {
		return false
	}
	return true

}
