package auth

import (
	"github.com/golang-jwt/jwt/v4"
)

// TokenClaims represents the authorization claims transmitted via a JWT.
type TokenClaims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}
