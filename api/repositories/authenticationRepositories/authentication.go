package authenticationRepositories

import (
	"divar.ir/internal/email"
	"divar.ir/internal/jwt"
	"divar.ir/internal/mongodb"
	"divar.ir/internal/password"
	"divar.ir/internal/redis"
	"divar.ir/internal/sessions_manager"
	"divar.ir/schema"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Register(ctx *gin.Context) (map[string]interface{}, error) {
	var registerInfo schema.UserInfo
	err := ctx.ShouldBind(&registerInfo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": false})
		return nil, err
	}
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	var existUser schema.UserInfo
	err = handler.FindOne("users", bson.M{"email": registerInfo.Email}, bson.M{}).Decode(&existUser)
	if !errors.Is(err, mongo.ErrNoDocuments) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email already exist", "status": false})
		return nil, errors.New("email already exist")
	}
	err = handler.FindOne("users", bson.M{"loginInfo.username": registerInfo.LoginInfo.Username}, bson.M{}).Decode(&existUser)
	if !errors.Is(err, mongo.ErrNoDocuments) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username already exist", "status": false})
		return nil, errors.New("username already exist")
	}
	hashedPassword, err := password.Hash(registerInfo.LoginInfo.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": "Error while hashing password"})
		return nil, err
	}
	registerInfo.LoginInfo.Password = hashedPassword
	randomNumber := rand.Intn(90000) + 10000
	registerInfo.LoginInfo.VerificationCode = strconv.Itoa(randomNumber)

	registerInfo.Id = primitive.NewObjectID()
	registerInfo.RegisterTime = time.Now()
	registerInfo.ActiveSessions = make([]schema.Session, 0)
	registerInfo.InactiveSessions = make([]schema.Session, 0)
	registerInfo.Status = true

	registerInfo.Image = "http://divar.test/userImages/defult_user_profile.png"
	result, err := handler.Insert("users", registerInfo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": "Error Inserting Database"})
		return nil, err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := email.SendMail(&wg, registerInfo.Email, "Verify Your Account", "Your Verification Code is "+strconv.Itoa(randomNumber))
		if err != nil {
			fmt.Println(err)
		}
	}()
	var session schema.Session
	session, err = jwt.Login(ctx, result.InsertedID.(primitive.ObjectID).Hex(), false)
	if err != nil {
		return nil, err
	}

	err = sessions_manager.OpenSession(result.InsertedID.(primitive.ObjectID), session)
	if err != nil {
		return nil, err
	}
	return gin.H{"status": true, "verify-token": session.AccessToken, "message": "User created successfully! Please check your email address"}, nil
}
func Login(ctx *gin.Context) (map[string]interface{}, error) {
	var user schema.Login
	var result schema.UserInfo
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	err = handler.FindOne("users", bson.M{"loginInfo.username": user.Username}, bson.M{}).Decode(&result)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": false, "error": "Invalid email or password"})
		return nil, err
	}

	if !password.Check(user.Password, result.LoginInfo.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": false, "error": "Invalid email or password"})
		return nil, err
	}
	if !result.Status {
		if result.LoginInfo.IsLocked {
			return nil, errors.New("your account is locked")
		}
		if result.LoginInfo.IsBanned {
			return nil, errors.New("your account is Banned")
		} else {
			return nil, errors.New("your account maybe deleted")
		}
	}
	var session schema.Session
	if result.LoginInfo.IsVerified {
		session, err = jwt.Login(ctx, result.Id.Hex(), true)
		if err != nil {
			return nil, err
		}
	} else {
		session, err = jwt.Login(ctx, result.Id.Hex(), false)
		if err != nil {
			return nil, err
		}
	}

	err = sessions_manager.OpenSession(result.Id, session)
	if err != nil {
		return nil, err
	}
	if result.LoginInfo.IsVerified {
		return gin.H{"status": true, "access-token": session.AccessToken, "refresh-token": session.RefreshToken}, nil
	} else {
		return gin.H{"status": true, "verify-token": session.AccessToken, "msg": "please check your email for get verification code"}, nil
	}
}
func Verify(ctx *gin.Context) (gin.H, error) {
	verifyCode, exist := ctx.GetPostForm("verifyCode")
	if !exist {
		return nil, errors.New("verify code missing")
	}
	verifyToken := ctx.GetHeader("Authorization")
	if verifyToken == "" {
		return nil, errors.New("missing Authorization header")
	}
	verifyToken = strings.TrimPrefix(verifyToken, "Bearer ")
	userId, sessionId, redisId, tokenType, err := jwt.Validate(verifyToken)
	if err != nil || tokenType != "verify-token" {
		return nil, errors.New("invalid verify token")
	}
	rid, err := redis.Get(sessionId + "-rid")
	if err != nil || rid != redisId {
		return nil, errors.New("invalid verify token")
	}
	var result schema.UserInfo
	userIdObj, _ := primitive.ObjectIDFromHex(userId)
	sessionIdObj, _ := primitive.ObjectIDFromHex(sessionId)
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return nil, errors.New("verification failed")
	}
	err = handler.FindOne("users", bson.M{"_id": userIdObj}, bson.M{"loginInfo": 1, "status": 1}).Decode(&result)
	if err != nil {
		return nil, errors.New("verification failed")
	}
	if !result.Status {
		if result.LoginInfo.IsLocked {
			return nil, errors.New("your account is locked")
		}
		if result.LoginInfo.IsBanned {
			return nil, errors.New("your account is Banned")
		} else {
			return nil, errors.New("your account maybe deleted")
		}
	}
	if result.LoginInfo.VerificationCode != verifyCode {
		return nil, errors.New("invalid verify code")
	}
	_, err = handler.Update("users", bson.M{"_id": userIdObj}, bson.M{"$set": bson.M{"loginInfo.isVerified": true}})
	if err != nil {
		return nil, errors.New("verification failed")
	}
	err = sessions_manager.CloseSession(userIdObj, sessionIdObj)
	if err != nil {
		return nil, errors.New("cannot generate new token")
	}
	session, err := jwt.Login(ctx, userId, true)
	if err != nil {
		return nil, errors.New("cannot generate new token")
	}
	err = sessions_manager.OpenSession(userIdObj, session)
	if err != nil {
		return nil, errors.New("cannot generate new token")
	}
	return gin.H{"status": true, "access-token": session.AccessToken, "refresh-token": session.RefreshToken}, nil
}
func Refresh(ctx *gin.Context) (map[string]interface{}, error) {
	refreshToken := ctx.GetHeader("Authorization")
	refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")
	userId, sessionId, _, tokenType, err := jwt.Validate(refreshToken)
	if err != nil {
		return nil, err
	}
	if tokenType != "refresh-token" {
		return nil, errors.New("invalid refresh token")
	}
	userIdObj, _ := primitive.ObjectIDFromHex(userId)
	sessionIdObj, _ := primitive.ObjectIDFromHex(sessionId)
	sessionInfo, err := sessions_manager.GetSession(userIdObj, sessionIdObj)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	expAccessToken, _ := strconv.ParseFloat(os.Getenv("JWT_EXP"), 64)
	if time.Now().Sub(sessionInfo.LastActivity).Hours() < expAccessToken {
		return nil, errors.New("token not expired")
	}
	err = sessions_manager.CloseSession(userIdObj, sessionInfo.ID)
	if err != nil {
		return nil, err
	}

	openTime := sessionInfo.OpenTime
	sessionInfo, _ = jwt.Login(ctx, userId, true)
	sessionInfo.ID = sessionIdObj
	sessionInfo.OpenTime = openTime
	err = sessions_manager.OpenSession(userIdObj, sessionInfo)
	if err != nil {
		return nil, err
	}

	return gin.H{"status": true, "access-token": sessionInfo.AccessToken, "refresh-token": sessionInfo.RefreshToken}, nil
}
func Logout(ctx *gin.Context) error {
	refreshToken := ctx.GetHeader("Authorization")
	refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")
	userId, sessionId, redisId, tokenType, err := jwt.Validate(refreshToken)
	if err != nil || tokenType != "access-token" {
		return errors.New("invalid access token")
	}
	rid, err := redis.Get(sessionId + "-rid")
	if err != nil || rid != redisId {
		return errors.New("invalid access token")
	}

	userIdObj, _ := primitive.ObjectIDFromHex(userId)
	sessionIdObj, _ := primitive.ObjectIDFromHex(sessionId)
	session, err := sessions_manager.GetSession(userIdObj, sessionIdObj)
	if err != nil {
		return errors.New("invalid access token")
	}
	session.LastActivity = time.Now()
	err = sessions_manager.CloseSession(userIdObj, sessionIdObj)
	if err != nil {
		return errors.New("logout failed")
	}
	err = sessions_manager.AddToInactiveSessions(userIdObj, session)
	if err != nil {
		return err
	}
	err = redis.Del(sessionId + "-rid")
	if err != nil {
		return errors.New("logout failed")
	}
	return nil
}
