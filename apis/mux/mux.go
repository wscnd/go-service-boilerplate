// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"net/http"
	"os"

	"github.com/wscnd/go-service-boilerplate/apis/middleware"
	"github.com/wscnd/go-service-boilerplate/libs/logger"
	"github.com/wscnd/go-service-boilerplate/libs/web"
)

type RouteAdder interface {
	Add(app *web.App, cfg Config)
}

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build    string
	Shutdown chan os.Signal
	Log      *logger.Logger
}

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(cfg Config, routeAdder RouteAdder) http.Handler {
	app := web.NewApp(cfg.Shutdown, middleware.Logger)

	routeAdder.Add(app, cfg)

	return app
}
