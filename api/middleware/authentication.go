package middleware

import (
	"divar.ir/internal/jwt"
	"divar.ir/internal/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func UserAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Authorization")
		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Required"})
			c.Abort()
			return
		}
		accessToken = strings.TrimPrefix(accessToken, "Bearer ")
		userId, sessionId, redisId, tokenType, err := jwt.Validate(accessToken)
		if err != nil || tokenType != "access-token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Access Token"})
			c.Abort()
			return
		}
		rid, err := redis.Get(sessionId + "-user-rid")
		if err != nil || rid != redisId {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token Not Authorized"})
			c.Abort()
			return
		}

		c.Set("userId", userId)
		c.Set("sessionId", sessionId)
		c.Next()

	}
}
func AdminAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Authorization")
		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Required"})
			c.Abort()
			return
		}
		accessToken = strings.TrimPrefix(accessToken, "Bearer ")
		userId, sessionId, redisId, tokenType, err := jwt.Validate(accessToken)
		if err != nil || tokenType != "access-token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Access Token"})
			c.Abort()
			return
		}
		rid, err := redis.Get(sessionId + "-admin-rid")
		if err != nil || rid != redisId {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token Not Authorized"})
			c.Abort()
			return
		}

		c.Set("userId", userId)
		c.Set("sessionId", sessionId)
		c.Next()

	}
}
