package types

import "github.com/golang-jwt/jwt/v5"

type JwtRefreshTokenClaims struct {
	TokenFamily string `json:"tokenFamily"`
	jwt.RegisteredClaims
}
