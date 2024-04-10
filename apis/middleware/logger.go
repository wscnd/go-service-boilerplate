package middleware

import (
	"context"
	"net/http"

	"github.com/wscnd/go-service-boilerplate/libs/web"
)

func Logger(handler web.Handler) web.Handler {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		err := handler(ctx, w, r)

		return err
	}
	return h
}
