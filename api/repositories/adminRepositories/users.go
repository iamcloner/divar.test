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
func PromoteToAdmin(userId primitive.ObjectID) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("server error")
	}
	_, err = handler.Update("users", bson.M{"_id": userId}, bson.M{"$set": bson.M{"loginInfo.isAdmin": true}})
	if err != nil {
		return errors.New("failed to promote to admin")
	}
	return nil
}
func LockUser(userId primitive.ObjectID) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("server error")
	}
	_, err = handler.Update("users", bson.M{"_id": userId}, bson.M{"$set": bson.M{"loginInfo.isLocked": true}})
	if err != nil {
		return errors.New("failed to lock user")
	}
	return nil
}
func UnlockUser(userId primitive.ObjectID) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("server error")
	}
	_, err = handler.Update("users", bson.M{"_id": userId}, bson.M{"$set": bson.M{"loginInfo.isLocked": false}})
	if err != nil {
		return errors.New("failed to lock user")
	}
	return nil
}
func BanUser(userId primitive.ObjectID) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("server error")
	}
	_, err = handler.Update("posts", bson.M{"userId": userId, "status": "Confirmed"}, bson.M{"$set": bson.M{"status": "Ban/Confirmed"}})
	_, err = handler.Update("posts", bson.M{"userId": userId, "status": "Pending"}, bson.M{"$set": bson.M{"status": "Ban/Pending"}})
	_, err = handler.Update("users", bson.M{"_id": userId}, bson.M{"$set": bson.M{"loginInfo.isBanned": true}})
	if err != nil {
		return errors.New("failed to ban user")
	}
	return nil
}
func UnbanUser(userId primitive.ObjectID) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("server error")
	}
	_, err = handler.Update("posts", bson.M{"userId": userId, "status": "Ban/Confirmed"}, bson.M{"$set": bson.M{"status": "Confirmed"}})
	_, err = handler.Update("posts", bson.M{"userId": userId, "status": "Ban/Pending"}, bson.M{"$set": bson.M{"status": "Pending"}})
	_, err = handler.Update("users", bson.M{"_id": userId}, bson.M{"$set": bson.M{"loginInfo.isBanned": false}})
	if err != nil {
		return errors.New("failed to lock user")
	}
	return nil
}
func TerminateUserSession(userId primitive.ObjectID, sessionId primitive.ObjectID) error {
	session, admin, err := sessions_manager.GetSession(userId, sessionId)
	if err != nil {
		return err
	}
	if admin {
		return errors.New("cannot terminate admin session")
	}
	err = sessions_manager.CloseSession(userId, sessionId)
	if err != nil {
		return err
	}
	err = sessions_manager.AddToInactiveSessions(userId, session)
	if err != nil {
		return err
	}
	return nil
}
func TerminateAllSessionUser(userId primitive.ObjectID) error {
	sessions, err := sessions_manager.GetSessions(userId)
	if err != nil {
		return err
	}
	for _, session := range sessions {
		err = sessions_manager.CloseSession(userId, session.ID)
		if err != nil {
			return err
		}
		err = sessions_manager.AddToInactiveSessions(userId, session)
		if err != nil {
			return err
		}
	}
	return nil
}
