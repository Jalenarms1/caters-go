package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(uid string) (string, error) {
	signingKey := os.Getenv("JWT_SECRET")

	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": uid,
		"exp": time.Now().Add(((24 * 365) * time.Hour)).Unix(),
	})

	token, err := claim.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return token, nil

}
