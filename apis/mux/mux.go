// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"net/http"
	"os"

	"github.com/wscnd/go-service-boilerplate/libs/logger"
)

type RouteAdder interface {
	Add(mux *http.ServeMux, cfg Config)
}

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build    string
	Shutdown chan os.Signal
	Log      *logger.Logger
}

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(cfg Config, routeAdder RouteAdder) *http.ServeMux {
	mux := http.NewServeMux()
	routeAdder.Add(mux, cfg)
	return mux
}
