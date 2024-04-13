package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/wscnd/go-service-boilerplate/apis/errs"
	"github.com/wscnd/go-service-boilerplate/libs/logger"
	"github.com/wscnd/go-service-boilerplate/libs/web"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *logger.Logger) web.MiddlewareHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			err := handler(ctx, w, r)
			if err == nil {
				return nil
			}

			log.Error(ctx, "message", "msg", err)

			var er errs.Error
			var status int

			// Construct the error document that we are going to send out.
			switch {
			// Trusted error, we know what was sent.
			case errs.IsError(err):
				er = errs.GetError(err)
				status = er.Status

      // Some other error that we don't know what to do (yet).
			default:
				er = errs.Error{
					Err:    errors.New("unknown"),
					Status: http.StatusInternalServerError,
				}
				status = http.StatusInternalServerError
			}

			if err := web.RespondJSON(ctx, w, er, status); err != nil {
				return err
			}

			// If we receive the shutdown err we need to return it
			// back to the base handler to shut down the service.
			if web.IsShutdown(err) {
				return err
			}

			return nil
		}

		return h
	}

	return m
}
