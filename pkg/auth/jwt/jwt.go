package jwt

import (
	"errors"
	"github.com/google/uuid"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(userID uuid.UUID, username string, tokenType string, exp time.Duration, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(exp).Unix(),
		"type":     tokenType,
	})

	return token.SignedString([]byte(secretKey))
}

func ValidateToken(tokenString string, secretKey string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secretKey), nil
	})
}
