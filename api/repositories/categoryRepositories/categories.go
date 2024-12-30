package categoryRepositories

import (
	"divar.ir/internal/mongodb"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

func GetCategories(code int) (interface{}, error) {
	handler, err := mongodb.GetMongoDBHandler()

	if err != nil {
		return nil, errors.New("operation Failed")
	}

	if code == 0 {
		result, err := handler.FindMany("categories", bson.M{}, bson.M{"_id": 0, "subs.filters": 0})
		if err != nil {
			return nil, errors.New("operation Failed")
		}
		return result, nil
	} else {
		sub := code % 100
		cat := code / 100
		if cat < 1 {
			var tmp bson.M
			err := handler.FindOne("categories", bson.M{"code": code}, bson.M{"_id": 0, "subs.filters": 0}).Decode(&tmp)
			if err != nil {

				return nil, errors.New("operation Failed")
			}
			return tmp, nil

		} else {

			var tmp bson.M
			err := handler.FindOne("categories", bson.M{"code": cat, "subs.code": sub}, bson.M{"subs.$": 1, "_id": 0}).Decode(&tmp)
			if err != nil {

				return nil, errors.New("operation Failed")
			}

			tmp = (tmp["subs"].(bson.A)[0]).(bson.M)
			tmp["code"] = cat*100 + sub
			return tmp, nil
		}
	}

}
