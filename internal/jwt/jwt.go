package jwt

import (
	"divar.ir/internal/password"
	"divar.ir/internal/redis"
	"divar.ir/schema"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"strconv"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func Generate(userId string, sessionId string, redisId string, tokenType string, exp int) (string, error) {
	var claims = jwt.MapClaims{}
	if tokenType != "refresh-token" {
		claims = jwt.MapClaims{
			"uid":  userId,
			"sid":  sessionId,
			"rid":  redisId,
			"type": tokenType,
			"exp":  time.Now().Add(time.Minute * time.Duration(exp)).Unix(),
		}
	} else {
		claims = jwt.MapClaims{
			"uid":  userId,
			"sid":  sessionId,
			"type": tokenType,
			"exp":  time.Now().Add(time.Hour * time.Duration(exp)).Unix(),
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func Validate(tokenString string) (uid string, sid string, rid string, tokenType string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil {
		return "", "", "", "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if tokenType, exists := claims["type"].(string); exists {
			rid, exists := claims["rid"].(string)
			sid, exists := claims["sid"].(string)
			if !exists {
				return "", "", "", "", jwt.ErrInvalidKey
			}
			if uid, exists := claims["uid"].(string); exists {
				return uid, sid, rid, tokenType, nil
			} else {
				return "", "", "", "", jwt.ErrInvalidKey
			}
		} else {
			return "", "", "", "", jwt.ErrInvalidKey
		}

	} else {
		return "", "", "", "", jwt.ErrTokenUnverifiable
	}
}
func Login(ctx *gin.Context, id string, verfied bool) (schema.Session, error) {
	var session schema.Session
	expAccessToken, _ := strconv.Atoi(os.Getenv("JWT_EXP"))
	expRefreshToken, _ := strconv.Atoi(os.Getenv("JWT_REF_EXP"))
	session.ID = primitive.NewObjectID()
	rid, err := password.GenerateRedisId()
	if err != nil {
		return session, err
	}
	var accessToken, refreshToken string
	if verfied {
		accessToken, err = Generate(id, session.ID.Hex(), rid, "access-token", expAccessToken)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": "Could not generate token"})
			return session, err
		}
		refreshToken, err = Generate(id, session.ID.Hex(), "", "refresh-token", expRefreshToken)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": "Could not generate token"})
			return session, err
		}
	} else {
		accessToken, err = Generate(id, session.ID.Hex(), rid, "verify-token", expAccessToken)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": "Could not generate token"})
			return session, err
		}
		refreshToken = ""
	}

	ip := ctx.ClientIP()
	platform := ctx.GetHeader("user-agent")

	session.IP = ip
	session.Platform = platform
	session.OpenTime = time.Now()
	session.LastActivity = time.Now()
	session.AccessToken = accessToken
	session.RefreshToken = refreshToken
	err = redis.Set(session.ID.Hex()+"-rid", rid, time.Duration(expAccessToken)*time.Minute)
	if err != nil {
		return schema.Session{}, err
	}
	return session, nil
}
