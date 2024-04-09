// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/wscnd/go-service-boilerplate/foundation/logger"
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

	h := func(w http.ResponseWriter, r *http.Request) {
		status := struct {
			Status string
		}{
			Status: "ok",
		}
    w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&status)
	}

	mux.HandleFunc("GET /", h)
	return mux
}
