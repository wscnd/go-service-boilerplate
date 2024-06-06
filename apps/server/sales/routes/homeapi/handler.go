package homeapi

import (
	"github.com/wscnd/go-service-boilerplate/apis/auth"
	"github.com/wscnd/go-service-boilerplate/apis/middleware"
	"github.com/wscnd/go-service-boilerplate/apis/mux"
	"github.com/wscnd/go-service-boilerplate/libs/web"
)

func Routes(app *web.App, cfg mux.Config) {
	app.Handle("GET /", homeHandler)

	// TESTING ERROR HANDLING
	app.Handle("GET /error", handlerWithError)

	// TESTING PANICS HANDLING
	app.Handle("GET /panics", handlerWithPanic)

	// TESTING AUTH HANDLING
	authen := middleware.Authenticate(cfg.Auth)
	applyAdminRule := middleware.Authorize(cfg.Auth, auth.RuleAdminOnly)
	app.Handle("GET /authn", handlerWithAuth, authen, applyAdminRule)
}
