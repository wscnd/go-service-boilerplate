package homeapi

import (
	"github.com/wscnd/go-service-boilerplate/libs/web"
)

func Routes(app *web.App) {
	app.HandleFunc("GET /", homeHandler)
}