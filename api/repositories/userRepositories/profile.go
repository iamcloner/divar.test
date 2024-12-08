package userRepositories

import (
	"divar.ir/internal/mongodb"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

func GetProfile(userId string) (bson.M, error) {
	userIdObj, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errors.New("access token is not valid")
	}
	var result bson.M
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return nil, errors.New("internal server error")
	}
	err = handler.FindOne("users", bson.M{"_id": userIdObj}, bson.M{"name": 1, "email": 1, "country": 1, "image": 1, "birthday": 1}).Decode(&result)
	if err != nil {
		return nil, errors.New("user not found")
	}
	email := result["email"].(string)
	mailPostfix := email[strings.LastIndex(email, "@"):]
	result["email"] = email[0:3] + "*********" + mailPostfix
	return result, nil
}
