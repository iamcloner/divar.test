package adminRepositories

import (
	"divar.ir/internal/mongodb"
	"divar.ir/schema"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUsers() ([]bson.M, error) {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return nil, errors.New("server error")
	}
	result, err := handler.FindMany("users", bson.M{}, bson.M{"loginInfo": 0, "activeSessions": 0, "inactiveSessions": 0})
	if err != nil {
		return nil, errors.New("server error")
	}
	return result, nil
}
func GetUser(userId primitive.ObjectID) (schema.UserInfo, error) {
	var result schema.UserInfo
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return result, errors.New("server error")
	}
	err = handler.FindOne("users", bson.M{"_id": userId}, bson.M{"loginInfo": 0, "activeSessions": 0, "inactiveSessions": 0}).Decode(&result)
	if err != nil {
		return result, errors.New("server error")
	}
	return result, nil
}
func GetUserLoginInfo(userId primitive.ObjectID) (schema.Login, error) {
	var result schema.UserInfo
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return result.LoginInfo, errors.New("server error")
	}
	err = handler.FindOne("users", bson.M{"_id": userId}, bson.M{"loginInfo": 1}).Decode(&result)
	result.LoginInfo.Password = "*****"
	if err != nil {
		return result.LoginInfo, err
	}
	return result.LoginInfo, nil
}
