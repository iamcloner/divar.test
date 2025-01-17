package sessions_manager

import (
	"divar.ir/internal/mongodb"
	"divar.ir/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"strconv"
	"time"
)

func OpenSession(userId primitive.ObjectID, sessionInfo schema.Session) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	_, err = handler.Update("users", bson.M{"_id": userId}, bson.M{"$addToSet": bson.M{"activeSessions": sessionInfo}})
	if err != nil {
		return err
	}
	return err
}

func GetSessions(userId primitive.ObjectID) ([]schema.Session, error) {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	var result schema.UserInfo
	err = handler.FindOne("users", bson.M{"_id": userId}, bson.M{"activeSessions": 1}).Decode(&result)
	if err != nil {
		return nil, err
	}
	sessions, err := updateSessions(userId, result.ActiveSessions)
	if err != nil {
		return nil, err
	}
	return sessions, err
}
func CloseSession(userId primitive.ObjectID, sessionId primitive.ObjectID) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	_, err = handler.Update("users", bson.M{"_id": userId}, bson.M{"$pull": bson.M{"activeSessions": bson.M{"_id": sessionId}}})
	return err
}
func updateSessions(userId primitive.ObjectID, sessions []schema.Session) ([]schema.Session, error) {
	var newSession []schema.Session
	expRefreshToken, _ := strconv.ParseFloat(os.Getenv("JWT_REF_EXP"), 64)
	for _, session := range sessions {
		openTime := session.OpenTime
		estTime := time.Now().Sub(openTime).Hours()
		if estTime >= expRefreshToken {
			err := CloseSession(userId, session.ID)
			if err != nil {
				return nil, err
			}
			err = AddToInactiveSessions(userId, session)
			if err != nil {
				return nil, err
			}
		} else {
			newSession = append(newSession, session)
		}

	}
	return newSession, nil
}
func GetSession(userId primitive.ObjectID, sessionId primitive.ObjectID) (schema.Session, bool, error) {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return schema.Session{}, false, err
	}
	var result schema.UserInfo
	err = handler.FindOne("users", bson.M{"_id": userId, "activeSessions": bson.M{"$elemMatch": bson.M{"_id": sessionId}}}, bson.M{"activeSessions.$": 1, "status": 1, "loginInfo.isLocked": 1, "loginInfo.isBanned": 1, "loginInfo.isAdmin": 1}).Decode(&result)
	if err != nil {
		return schema.Session{}, false, err
	}

	return result.ActiveSessions[0], result.LoginInfo.IsAdmin, err
}
func AddToInactiveSessions(userId primitive.ObjectID, sessionInfo schema.Session) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	_, err = handler.Update("users", bson.M{"_id": userId}, bson.M{"$addToSet": bson.M{"inactiveSessions": sessionInfo}})
	if err != nil {
		return err
	}
	return err
}
