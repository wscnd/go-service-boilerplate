// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"net/http"
	"os"

	"github.com/wscnd/go-service-boilerplate/apps/server/sales/routes"
	"github.com/wscnd/go-service-boilerplate/libs/logger"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build    string
	Shutdown chan os.Signal
	Log      *logger.Logger
}

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(cfg Config) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", routes.Home)
	return mux
}
