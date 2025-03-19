package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
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

func GenerateRandomUrlSlug() string {
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"

	slug := make([]byte, 16)
	for i := range 16 {
		slug[i] = charset[rand.Intn(len(charset))]
	}

	return string(slug)
}

func SaveImage(path, dataUrl string) error {
	parts := strings.Split(dataUrl, ",")
	fmt.Println(parts)
	if len(parts) < 2 {
		return errors.New("invalid image")
	}

	image, err := base64.StdEncoding.DecodeString(parts[1])

	if err != nil {
		return err
	}

	return os.WriteFile(path, image, 0644)
}
