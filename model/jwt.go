package model

import "github.com/golang-jwt/jwt/v4"

type JwtCustomClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}
