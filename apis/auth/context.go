package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type ctxKey int

const (
	claimKey ctxKey = iota + 1
	userIDKey
)

func SetClaims(ctx context.Context, claims TokenClaims) context.Context {
	return context.WithValue(ctx, claimKey, claims)
}

func GetClaims(ctx context.Context) TokenClaims {
	v, ok := ctx.Value(claimKey).(TokenClaims)
	if !ok {
		return TokenClaims{}
	}
	return v
}

func SetUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID returns the claims from the context.
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("user id not found in context")
	}

	return v, nil
}
