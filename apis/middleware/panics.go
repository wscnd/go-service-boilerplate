package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/wscnd/go-service-boilerplate/apis/metrics"
	"github.com/wscnd/go-service-boilerplate/libs/web"
)

// Panics recovers from panics and converts the panic to an error so it is
// reported in Metrics and handled in Errors.
func Panics() web.MiddlewareHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			// Needs to use named return since defer statements are kind of in a limbo space.
			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))

					metrics.AddPanics(ctx)
				}
			}()

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
