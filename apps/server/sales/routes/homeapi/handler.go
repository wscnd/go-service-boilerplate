package homeapi

import (
	"github.com/wscnd/go-service-boilerplate/apis/mux"
	"github.com/wscnd/go-service-boilerplate/libs/web"
)

func Routes(app *web.App, cfg mux.Config) {
	app.Handle("GET /", homeHandler)

	// TESTING ERROR HANDLING
	app.Handle("GET /error", handlerWithError)

	// TESTING PANICS HANDLING
	app.Handle("GET /panics", handlerWithPanic)

}
