package listPicker

import (
	"divar.ir/internal/mongodb"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

func GetList(name string) ([]string, error) {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return nil, errors.New("operation Failed")
	}
	var result bson.M
	err = handler.FindOne("lists", bson.M{"name": name}, bson.M{"_id": 0, "name": 0, "values": 1}).Decode(&result)
	if err != nil {
		return nil, errors.New("operation Failed")
	}
	list := result["values"].([]string)
	return list, nil
}
func AddList(name string, list []string) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("operation Failed")
	}
	_, err = handler.Insert("lists", bson.M{"name": name, "values": list})
	if err != nil {
		return errors.New("operation Failed")
	}
	return nil
}
func CheckItem(name string, item string) (bool, error) {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return false, errors.New("operation Failed")
	}
	var result bson.M
	err = handler.FindOne("lists", bson.M{"name": name, "values": item}, bson.M{}).Decode(&result)
	if err != nil {
		return false, nil
	}
	return true, nil
}
func ExistList(name string) (bool, error) {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return false, errors.New("operation Failed")
	}
	var result bson.M
	err = handler.FindOne("lists", bson.M{"name": name}, bson.M{}).Decode(&result)
	if err != nil {
		return false, errors.New("operation Failed")
	}
	return true, nil
}
func AddItem(name string, item string) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("operation Failed")
	}
	_, err = handler.Update("lists", bson.M{"name": name}, bson.M{"$addToSet": bson.M{"values": item}})
	if err != nil {
		return errors.New("operation Failed")
	}
	return nil
}
func RemoveItem(name string, item string) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("operation Failed")
	}
	_, err = handler.Update("lists", bson.M{"name": name}, bson.M{"$pull": bson.M{"values": item}})
	if err != nil {
		return errors.New("operation Failed")
	}
	return nil
}
func RemoveList(name string) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("operation Failed")
	}
	_, err = handler.Delete("lists", bson.M{"name": name})
	if err != nil {
		return errors.New("operation Failed")
	}
	return nil
}
