package middleware

import (
	"context"
	"net/http"

	"github.com/wscnd/go-service-boilerplate/apis/metrics"
	"github.com/wscnd/go-service-boilerplate/libs/web"
)

// Metrics updates metrics counters.
func Metrics() web.MiddlewareHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx = metrics.Set(ctx)

			err := handler(ctx, w, r)
			if err != nil {
				metrics.AddErrors(ctx)
			}

			n := metrics.AddRequests(ctx)

			// 1_000 is a arbitraty number
			if n%1_000 == 0 {
				metrics.AddGoroutines(ctx)
			}

			return err
		}

		return h
	}

	return m
}
