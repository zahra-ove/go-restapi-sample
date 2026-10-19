package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secret = "supersecet"

func GenerateToken(email string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func VerifyToken(token string) error {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return errors.New("Could not parse token")
	}

	if !parsedToken.Valid {
		return errors.New("Invalid token")
	}

	return nil
}
