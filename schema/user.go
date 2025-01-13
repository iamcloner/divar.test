package schema

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Session struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	AccessToken  string             `json:"-" form:"-" bson:"accessToken"`
	RefreshToken string             `json:"-" form:"-" bson:"refreshToken"`
	IP           string             `json:"-" form:"-" bson:"ip"`
	OpenTime     time.Time          `json:"-" form:"-" bson:"openTime"`
	Platform     string             `json:"-" form:"-" bson:"platform"`
	LastActivity time.Time          `json:"-" form:"-" bson:"lastActivity"`
}

type Login struct {
	Username         string `json:"username" form:"username" binding:"required,min=4,max=16" bson:"username"`
	Password         string `json:"password" form:"password" binding:"required,min=8,max=26" bson:"password"`
	IsAdmin          bool   `json:"isAdmin" form:"-" bson:"isAdmin"`
	IsVerified       bool   `json:"isVerified" form:"-" bson:"isVerified"`
	IsLocked         bool   `json:"isLocked" form:"-" bson:"isLocked"`
	IsBanned         bool   `json:"isBanned" form:"-" bson:"isBanned"`
	VerificationCode string `json:"verificationCode" form:"-" bson:"verificationCode"`
}

type UserInfo struct {
	Id               primitive.ObjectID `json:"-" form:"-" bson:"_id"`
	Status           bool               `json:"-" form:"-" bson:"status"`
	Name             string             `json:"name" form:"name" binding:"required,min=5,max=16" bson:"name"`
	Email            string             `json:"email" form:"email" binding:"required,email" bson:"email"`
	IsAdmin          bool               `json:"-" form:"-" bson:"isAdmin"`
	Birthday         time.Time          `json:"birthday" form:"birthday" binding:"required" bson:"birthday"`
	Address          string             `json:"address" form:"address" binding:"required,min=5,max=60" bson:"address"`
	Country          string             `json:"country" form:"country" binding:"required,min=3,max=30" bson:"country"`
	LoginInfo        Login              `json:"-" form:"login" binding:"required" bson:"loginInfo" bson:"loginInfo"`
	RegisterTime     time.Time          `json:"-" form:"-" bson:"registerTime"`
	Image            string             `json:"-" form:"-" bson:"image"`
	ActiveSessions   []Session          `json:"-" form:"-" bson:"activeSessions"`
	InactiveSessions []Session          `json:"-" form:"-" bson:"inactiveSessions"`
}
