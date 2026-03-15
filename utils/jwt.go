package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/BlitzStudio/blitzStudioAuth/types"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userId int64, tokenId string, timeNow time.Time) string {
	keyData := os.Getenv("JWT_ACCESS_TOKEN_KEY")
	if len(keyData) == 0 {
		logger.Fatal("Missing JWT_ACCESS_TOKEN_KEY")
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(keyData))
	if err != nil {
		logger.Fatal(err)
	}

	// doar pentru a genera accessTokens
	claims := jwt.RegisteredClaims{
		ID:        tokenId,
		Subject:   fmt.Sprintf("%d", userId),
		IssuedAt:  jwt.NewNumericDate(timeNow),
		NotBefore: jwt.NewNumericDate(timeNow),
		ExpiresAt: jwt.NewNumericDate(timeNow.Add(15 * time.Minute)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		logger.Fatal(err)
	}
	return signedToken
}

func GenerateRefreshToken(userId int64, tokenId string, tokenFamily string, timeNow time.Time) string {
	keyData := os.Getenv("JWT_REFRESH_TOKEN_KEY")
	if len(keyData) == 0 {
		logger.Fatal("Missing JWT_REFRESH_TOKEN_KEY")
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(keyData))

	if err != nil {
		logger.Fatal(err)
	}

	claims := types.JwtRefreshTokenClaims{
		TokenFamily: tokenFamily,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId,
			Subject:   fmt.Sprintf("%d", userId),
			IssuedAt:  jwt.NewNumericDate(timeNow),
			NotBefore: jwt.NewNumericDate(timeNow),
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(7 * 24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		logger.Fatal(err)
	}
	return signedToken
}

func ValidateRefreshToken(tokenString string) (*types.JwtRefreshTokenClaims, error) {
	keyData := os.Getenv("JWT_REFRESH_TOKEN_PUB_KEY")
	if len(keyData) == 0 {
		logger.Fatal("Missing JWT_REFRESH_TOKEN_PUB_KEY")
	}
	token, err := jwt.ParseWithClaims(tokenString, &types.JwtRefreshTokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			logger.Warn("Someone tried using a symmetric key")
			return nil, errors.New("Someone tried using a symmetric key")
		}
		pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(keyData))
		if err != nil {
			logger.Fatal(err)
		}
		return pubKey, nil
	})
	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*types.JwtRefreshTokenClaims); ok {
		return claims, nil
	}
	return nil, errors.New("Token validation error")
}
