package adminRepositories

import (
	"divar.ir/internal/mongodb"
	"divar.ir/internal/sessions_manager"
	"divar.ir/schema"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
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
func GetUserActiveSession(userId primitive.ObjectID) ([]schema.Session, error) {
	sessions, err := sessions_manager.GetSessions(userId)
	if err != nil {
		return []schema.Session{}, errors.New("failed to get active sessions")
	}
	return sessions, err
}
func GetUserInactiveSession(userId primitive.ObjectID) ([]schema.Session, error) {
	_, err := sessions_manager.GetSessions(userId)
	if err != nil {
		return []schema.Session{}, errors.New("failed to update")
	}
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	var result schema.UserInfo
	err = handler.FindOne("users", bson.M{"_id": userId}, bson.M{"inactiveSessions": 1}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.InactiveSessions, err
}
