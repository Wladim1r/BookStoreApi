package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	JWTSecretKey = "GIN JWT POSTGRESQL DOCKER"
	TokenExpity  = 24 * time.Hour
)

func GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(TokenExpity).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(JWTSecretKey))
}

func ParseToken(signedToken string) (*jwt.Token, error) {
	return jwt.Parse(signedToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(JWTSecretKey), nil
	})
}
