package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userId uint, email string, duration time.Duration) (string, error) {
	claim := &jwt.MapClaims{
		"userId": userId,
		"email":  email,
		"exp":    time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(jwtKey)
}

func ValidateJWT(signedToken string) (*jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)

	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}
