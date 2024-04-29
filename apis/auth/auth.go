package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/open-policy-agent/opa/rego"
	"github.com/wscnd/go-service-boilerplate/libs/logger"
)

// TokenClaims represents the authorization claims transmitted via a JWT.
type TokenClaims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

// KeyLookup declares a method set of behavior for looking up
// private and public keys for JWT use. The return could be a
// PEM encoded string or a JWS based key.
type KeyLookup interface {
	PrivateKey(keyID string) (key string, err error)
	PublicKey(keyID string) (key string, err error)
}

// Config represents information required to initialize auth.
type Config struct {
	Log       *logger.Logger
	KeyLookup KeyLookup
	Issuer    string
}

// Auth is used to authenticate clients. It can generate a token for a
// set of user claims and recreate the claims by parsing the token.
type Auth struct {
	keyLookup KeyLookup
	method    jwt.SigningMethod
	parser    *jwt.Parser
	issuer    string
}

// New creates an Auth to support authentication/authorization.
func New(cfg Config) (*Auth, error) {
	a := Auth{
		keyLookup: cfg.KeyLookup,
		method:    jwt.GetSigningMethod(jwt.SigningMethodRS256.Name),
		parser:    jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name})),
		issuer:    cfg.Issuer,
	}

	return &a, nil
}

// GenerateToken generates a signed JWT token string representing the user TokenClaims.
func (a *Auth) GenerateToken(keyID string, claims TokenClaims) (string, error) {
	token := jwt.NewWithClaims(a.method, claims)
	token.Header["keyID"] = keyID

	privateKeyPEM, err := a.keyLookup.PrivateKey(keyID)
	if err != nil {
		return "", fmt.Errorf("private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		return "", fmt.Errorf("parsing private pem: %w", err)
	}

	tokenStr, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return tokenStr, nil
}

// Authenticate processes the token to validate the sender's token is valid.
func (a *Auth) Authenticate(ctx context.Context, bearerToken string) (TokenClaims, error) {
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return TokenClaims{}, errors.New("expected authorization header format: Bearer <token>")
	}

	var claims TokenClaims
	token, _, err := a.parser.ParseUnverified(parts[1], &claims)
	if err != nil {
		return TokenClaims{}, fmt.Errorf("error parsing token: %w", err)
	}

	keyIDraw, exists := token.Header["keyID"]
	if !exists {
		return TokenClaims{}, fmt.Errorf("keyID missing from header: %w", err)
	}

	keyID, ok := keyIDraw.(string)
	if !ok {
		return TokenClaims{}, fmt.Errorf("keyID malformed: %w", err)
	}

	pemStr, err := a.keyLookup.PublicKey(keyID)
	if err != nil {
		return TokenClaims{}, fmt.Errorf("failed to fetch public key: %w", err)
	}

	opaInput := map[string]any{
		"Key":   pemStr,
		"Token": parts[1],
	}

	if err := a.opaPolicyEvaluation(ctx, opaAuthentication, RuleAuthenticate, opaInput); err != nil {
		return TokenClaims{}, fmt.Errorf("authentication failed : %w", err)
	}

	return claims, nil
}

// Authorize attempts to authorize the user with the provided input roles, if
// none of the input roles are within the user's claims, we return an error
// otherwise the user is authorized.
func (a *Auth) Authorize(ctx context.Context, claims TokenClaims, userID uuid.UUID, rule string) error {
	input := map[string]any{
		"Roles":   claims.Roles,
		"Subject": claims.Subject,
		"UserID":  userID,
	}

	if err := a.opaPolicyEvaluation(ctx, opaAuthorization, rule, input); err != nil {
		return fmt.Errorf("rego evaluation failed : %w", err)
	}

	return nil
}

func (a *Auth) opaPolicyEvaluation(ctx context.Context, opaPolicy string, rule string, input any) error {
	query := fmt.Sprintf("x = data.%s.%s", opaPackage, rule)

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
