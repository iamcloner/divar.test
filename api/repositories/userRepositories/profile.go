package userRepositories

import (
	"divar.ir/internal/email"
	"divar.ir/internal/mongodb"
	"divar.ir/internal/password"
	"divar.ir/internal/redis"
	"divar.ir/schema"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

func GetProfile(userId string) (schema.UserInfo, error) {
	userIdObj, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return schema.UserInfo{}, errors.New("access token is not valid")
	}
	var result schema.UserInfo
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return schema.UserInfo{}, errors.New("internal server error")
	}
	err = handler.FindOne("users", bson.M{"_id": userIdObj}, bson.M{"name": 1, "email": 1, "country": 1, "image": 1, "birthday": 1, "address": 1, "status": 1, "loginInfo.isLocked": 1, "loginInfo.isBanned": 1}).Decode(&result)
	if err != nil {
		return schema.UserInfo{}, errors.New("user not found")
	}
	if !result.Status {
		if result.LoginInfo.IsLocked {
			return schema.UserInfo{}, errors.New("your account is locked")
		}
		if result.LoginInfo.IsBanned {
			return schema.UserInfo{}, errors.New("your account is Banned")
		} else {
			return schema.UserInfo{}, errors.New("your account maybe deleted")
		}
	}
	mailPostfix := result.Email[strings.LastIndex(result.Email, "@"):]
	result.Email = result.Email[0:3] + "*********" + mailPostfix
	return result, nil
}
func UpdateProfileName(userId primitive.ObjectID, newName string) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("internal server error")
	}
	_, err = handler.Update("users", bson.M{"_id": userId}, bson.M{"$set": bson.M{"name": newName}})
	if err != nil {
		return errors.New("update Name failed")
	}
	return nil
}
func UpdateProfileEmail(userId primitive.ObjectID, oldEmail string) error {
	var result schema.UserInfo
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("internal server error")
	}
	err = handler.FindOne("users", bson.M{"_id": userId, "email": oldEmail}, bson.M{"email": 1, "status": 1, "loginInfo.isLocked": 1, "loginInfo.isBanned": 1}).Decode(&result)
	if err != nil {
		return errors.New("invalid Email")
	}
	verifyCode := strconv.Itoa(rand.Intn(90000) + 10000)
	err = redis.Set(userId.Hex()+"-changeEmail", verifyCode, 10*time.Minute)
	if err != nil {
		return errors.New("failed to update email")
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := email.SendMail(&wg, oldEmail, "Verify Your Account", "Your Verification Code is "+verifyCode)
		if err != nil {
			fmt.Println(err)
		}
	}()
	return nil
}
func ConfirmUpdateProfileEmail(userId primitive.ObjectID, verifyCode string, newEmail string) error {
	verifyCodeRedis, err := redis.Get(userId.Hex() + "-changeEmail")
	if err != nil {
		return errors.New("code has been expired")
	}
	if verifyCodeRedis != verifyCode {
		return errors.New("invalid code")
	}
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("internal server error")
	}
	_, err = handler.Update("users", bson.M{"_id": userId}, bson.M{"$set": bson.M{"email": newEmail}})
	if err != nil {
		return errors.New("update Email failed")
	}
	return nil
}
func ConfirmDeleteProfile(userId primitive.ObjectID, verifyCode string) error {
	verifyCodeRedis, err := redis.Get(userId.Hex() + "-deleteAccount")
	if err != nil {
		return errors.New("code has been expired")
	}
	if verifyCodeRedis != verifyCode {
		return errors.New("invalid code")
	}
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("internal server error")
	}
	_, err = handler.Update("users", bson.M{"_id": userId}, bson.M{"status": false})
	if err != nil {
		return errors.New("delete account failed")
	}
	return nil
}
func DeleteProfile(userId primitive.ObjectID) error {
	var result schema.UserInfo
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("internal server error")
	}
	err = handler.FindOne("users", bson.M{"_id": userId}, bson.M{"email": 1, "status": 1, "loginInfo.isLocked": 1, "loginInfo.isBanned": 1}).Decode(&result)
	if err != nil {
		return errors.New("invalid account")
	}
	verifyCode := strconv.Itoa(rand.Intn(90000) + 10000)
	err = redis.Set(userId.Hex()+"-deleteAccount", verifyCode, 10*time.Minute)
	if err != nil {
		return errors.New("failed to delete account")
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := email.SendMail(&wg, result.Email, "Verify Your Account", "Your Verification Code is "+verifyCode)
		if err != nil {
			fmt.Println(err)
		}
	}()
	return nil
}
func CheckOldPassword(userId primitive.ObjectID, oldPassword string) (bool, error) {
	var result schema.UserInfo
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return false, errors.New("internal server error")
	}
	err = handler.FindOne("users", bson.M{"_id": userId}, bson.M{"loginInfo.password": 1}).Decode(&result)
	if err != nil {
		return false, errors.New("invalid user id")
	}
	return password.Check(oldPassword, result.LoginInfo.Password), nil
}
func UpdatePassword(userId primitive.ObjectID, newPassword string) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("internal server error")
	}
	newPasswordHash, err := password.Hash(newPassword)
	if err != nil {
		return err
	}
	_, err = handler.Update("users", bson.M{"_id": userId}, bson.M{"$addToSet": bson.M{"password": newPasswordHash}})
	if err != nil {
		return errors.New("invalid user id")
	}
	return nil
}
