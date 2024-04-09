package handler

import (
	"net/http"
)

func Handler(mux *http.ServeMux) {
	mux.HandleFunc("GET /", homeHandler)
}