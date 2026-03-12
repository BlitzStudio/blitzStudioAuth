package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateToken(userId int64, timeNow time.Time, expirationTime time.Time) string {
	keyData := os.Getenv("JWT_PRIVATE_KEY")
	if len(keyData) == 0 {
		logger.Fatal("Missing JWT_PRIVATE_KEY")
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(keyData))
	if err != nil {
		logger.Fatal(err)
	}
	claims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userId),
		IssuedAt:  jwt.NewNumericDate(timeNow),
		NotBefore: jwt.NewNumericDate(timeNow),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		logger.Fatal(err)
	}
	return signedToken
}

func GenerateAccessToken(userId int64, timeNow time.Time) string {
	return generateToken(userId, timeNow, timeNow.Add(15*time.Minute))
}

func GenerateRefreshToken(userId int64, timeNow time.Time) string {
	return generateToken(userId, timeNow, timeNow.Add(7*24*time.Hour))
}
