package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/wscnd/go-service-boilerplate/apis/auth"
	"github.com/wscnd/go-service-boilerplate/apis/errs"

	"github.com/wscnd/go-service-boilerplate/libs/web"
)

var (
	ErrInvalidID       = errors.New("ID is not in its proper form")
	ErrUnauthenticated = errors.New("Unauthenticated")
)

// Authenticate validates a JWT from the `Authorization` header.
func Authenticate(a *auth.Auth) web.MiddlewareHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			claims, err := a.Authenticate(ctx, r.Header.Get("authorization"))
			if err != nil {
				return errs.New(ErrUnauthenticated, http.StatusUnauthorized)
			}

			if claims.Subject == "" {
				return errs.Newf(ErrUnauthenticated, http.StatusUnauthorized, "authorize: you are not authorized for that action, no claims")
			}

			subjectID, err := uuid.Parse(claims.Subject)
			if err != nil {
				return errs.New(fmt.Errorf("parsing subject: %w", err), http.StatusUnauthorized)
			}

			ctx = auth.SetUserID(ctx, subjectID)
			ctx = auth.SetClaims(ctx, claims)

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}

// Authorize executes the specified role and does not extract any domain data.
func Authorize(a *auth.Auth, rule string) web.MiddlewareHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			claims := auth.GetClaims(ctx)
			if err := a.Authorize(ctx, claims, uuid.UUID{}, rule); err != nil {
				return errs.Newf(ErrUnauthenticated, http.StatusUnauthorized, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, rule, err)
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
