package commands

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/wscnd/go-service-boilerplate/apis/auth"
	"github.com/wscnd/go-service-boilerplate/libs/logger"
)

func GenToken(log *logger.Logger) error {
	// Generating a token requires defining a set of tokenClaims. In this applications
	// case, we only care about defining the subject and the user in question and
	// the roles they have on the database. This token will expire in a year.
	//
	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expires
	// nbf (not before time): Time before which the JWT must not be accepted for processing
	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)
	tokenClaims := auth.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1234567890",
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)), /* a year */
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{"ADMIN"},
	}

	method := jwt.GetSigningMethod(jwt.SigningMethodRS256.Name)
	token := jwt.NewWithClaims(method, tokenClaims)
	token.Header["kid"] = "60877A3C-9AB6-4A50-9F27-B56D78229D92"

	tokenStr, err := token.SignedString(privateKey)
	if err != nil {
		return fmt.Errorf("signing token: %w", err)
	}

	fmt.Println("************")
	fmt.Println(tokenStr)
	fmt.Println("************")

	// Validating the jwt token
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}))
	var tokenClaims2 auth.TokenClaims
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return &privateKey.PublicKey, nil
	}

	tkn, err := parser.ParseWithClaims(tokenStr, &tokenClaims2, keyFunc)
	if err != nil {
		return fmt.Errorf("parsing token claims2: %w", err)
	}

	if !tkn.Valid {
		return fmt.Errorf("signature invalid: %w", err)
	}

	fmt.Println("************")
	fmt.Printf("%#v\n", tokenClaims2)
	fmt.Println("************")

	// TODO: validate with opa

	return nil
}
