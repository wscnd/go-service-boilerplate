package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/wscnd/go-service-boilerplate/apis/auth"
	"github.com/wscnd/go-service-boilerplate/libs/keystore"
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
			Subject:   "5cf37266-3473-4006-984f-9325122678b7",
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)), /* a year */
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{"user"},
	}

	ks := keystore.New()
	if err := ks.LoadRSAKeys(os.DirFS("zarf/keys/")); err != nil {
		return fmt.Errorf("reading keys: %w", err)
	}
	authCfg := auth.Config{
		Log:       log,
		KeyLookup: ks,
	}
	a, err := auth.New(authCfg)
	if err != nil {
		return fmt.Errorf("constructing auth: %w", err)
	}

	tokenStr, err := a.GenerateToken("60877A3C-9AB6-4A50-9F27-B56D78229D92", tokenClaims)
	if err != nil {
		return err
	}

	fmt.Printf("-----BEGIN TOKEN-----\n%s\n-----END TOKEN-----\n", tokenStr)

	claims, err := a.Authenticate(context.Background(), fmt.Sprintf("Bearer %s", tokenStr))
	if err != nil {
		return fmt.Errorf("parsing token claims: %w", err)
	}

	fmt.Println("************")
	fmt.Printf("%#v\n", claims)
	fmt.Println("************")

	return nil
}
