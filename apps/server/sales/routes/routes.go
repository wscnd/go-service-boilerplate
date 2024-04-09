package routes

import (
	"net/http"

	"github.com/wscnd/go-service-boilerplate/apis/mux"
	"github.com/wscnd/go-service-boilerplate/apps/server/sales/handler"
)


type Routes struct{}

func (Routes) Add(mux *http.ServeMux, cfg mux.Config) {
  handler.Handler(mux)
}