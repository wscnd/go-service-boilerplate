package routes

import (
	"github.com/wscnd/go-service-boilerplate/apis/mux"
	"github.com/wscnd/go-service-boilerplate/apps/server/sales/routes/homeapi"
	"github.com/wscnd/go-service-boilerplate/libs/web"
)

type Routes struct{}

func (Routes) Add(app *web.App, cfg mux.Config) {
	homeapi.Routes(app, cfg)
}
