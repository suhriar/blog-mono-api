package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/suhriar/blog-mono-api/config"
	"github.com/suhriar/blog-mono-api/model"
)

func GenerateJWT(userID int64, email, username string) (tokenString string, err error) {
	claims := &model.JwtCustomClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err = token.SignedString([]byte(config.AppConfig.Jwt.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (claims model.JwtCustomClaims, err error) {
	t, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.Jwt.Secret), nil
	})

	if err != nil {
		return claims, err
	}

	if !t.Valid {
		return claims, fmt.Errorf("invalid token")
	}

	return claims, nil
}
