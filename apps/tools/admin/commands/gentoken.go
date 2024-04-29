package commands

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/wscnd/go-service-boilerplate/apis/auth"

	"github.com/open-policy-agent/opa/rego"
)

func GenToken() error {
	privateKey, err := GenKey()
	if err != nil {
		return fmt.Errorf("generating key: %w", err)
	}
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
	tokenClaims := struct {
		jwt.RegisteredClaims
		Roles []string
	}{
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
	var tokenClaims2 struct {
		jwt.RegisteredClaims
		Roles []string
	}
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

	var tokenClaims3 struct {
		jwt.RegisteredClaims
		Roles []string
	}

	_, _, err = parser.ParseUnverified(tokenStr, &tokenClaims3)
	if err != nil {
		return fmt.Errorf("parsing token claims3: %w", err)
	}

	// Marshal the public key from the private key to PKIX.
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("marshaling public key: %w", err)
	}

	// Construct a PEM block for the public key.
	publicBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	var pemFile bytes.Buffer
	// Write the public key to the public key file.
	if err := pem.Encode(&pemFile, &publicBlock); err != nil {
		return fmt.Errorf("encoding to public file: %w", err)
	}

	opaInput := map[string]any{
		"Key":   pemFile.String(),
		"Token": tokenStr,
	}

	err = opaPolicyEvaluation(context.Background(), auth.OpaAuthentication, auth.RuleAuthenticate, opaInput)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	fmt.Println("************")
	fmt.Println("SIGNATURE VALIDATED WITH REGO")
	fmt.Println("************")

	return nil
}

func opaPolicyEvaluation(ctx context.Context, opaPolicy string, rule string, input any) error {
	query := fmt.Sprintf("x = data.%s.%s", auth.OpaPackage, rule)

	q, err := rego.New(
		rego.Query(query),
		rego.Module("policy.rego", opaPolicy),
	).PrepareForEval(ctx)
	if err != nil {
		return err
	}

	results, err := q.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	if len(results) == 0 {
		return errors.New("no results")
	}
	fmt.Printf("%#v\n", results)

	result, ok := results[0].Bindings["x"].(bool)
	if !ok || !result {
		return fmt.Errorf("bindings results[%v] ok[%v]", results, ok)
	}

	return nil
}
