package sessions_manager

import (
	"divar.ir/internal/mongodb"
	"divar.ir/schema"
	"errors"
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

func GetSessions(userId primitive.ObjectID) ([]interface{}, error) {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	var result bson.M
	err = handler.FindOne("users", bson.M{"_id": userId}, bson.M{"activeSessions": 1}).Decode(&result)
	if err != nil {
		return nil, err
	}
	sessions, ok := result["activeSessions"].(primitive.A)
	if !ok {
		println(ok)
	}
	sessions, err = updateSessions(userId, sessions)
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
func updateSessions(userId primitive.ObjectID, sessions []interface{}) (primitive.A, error) {
	var newSession []interface{}
	expRefreshToken, _ := strconv.ParseFloat(os.Getenv("JWT_REF_EXP"), 64)
	for _, session := range sessions {
		sessionData, _ := session.(bson.M)
		openTime := sessionData["lastActivity"].(primitive.DateTime)
		estTime := time.Now().Sub(openTime.Time()).Hours()
		if estTime >= expRefreshToken {
			err := CloseSession(userId, sessionData["_id"].(primitive.ObjectID))
			if err != nil {
				return nil, err
			}
			err = AddToInactiveSessions(userId, session.(schema.Session))
			if err != nil {
				return nil, err
			}
		} else {
			newSession = append(newSession, session)
		}

	}
	return newSession, nil
}
func GetSession(userId primitive.ObjectID, sessionId primitive.ObjectID) (schema.Session, error) {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return schema.Session{}, err
	}
	var result schema.UserInfo
	err = handler.FindOne("users", bson.M{"_id": userId, "activeSessions._id": sessionId}, bson.M{"activeSessions.$": 1, "status": 1, "loginInfo.isLocked": 1, "loginInfo.isBanned": 1}).Decode(&result)
	if err != nil {
		return schema.Session{}, err
	}
	if !result.Status {
		if result.LoginInfo.IsLocked {
			return schema.Session{}, errors.New("your account is locked")
		}
		if result.LoginInfo.IsBanned {
			return schema.Session{}, errors.New("your account is Banned")
		} else {
			return schema.Session{}, errors.New("your account maybe deleted")
		}
	}
	return result.ActiveSessions[0], err
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
