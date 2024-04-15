package homeapi

import (
	"github.com/wscnd/go-service-boilerplate/libs/web"
)

func Routes(app *web.App) {
	app.Handle("GET /", homeHandler)
	app.Handle("GET /error", handlerWithError)
	app.Handle("GET /panics", handlerWithPanic)
}
