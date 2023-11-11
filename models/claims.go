package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type AppClaims struct {
	UserId int64 `json:"user_id"`
	jwt.StandardClaims
}

