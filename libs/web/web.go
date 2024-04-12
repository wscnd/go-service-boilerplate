package web

import (
	"context"
	"net/http"
	"os"
	"syscall"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	*http.ServeMux
	shutdown chan os.Signal
	mws      []MiddlewareHandler
}

func NewApp(shutdown chan os.Signal, mws ...MiddlewareHandler) *App {
	mux := http.NewServeMux()
	return &App{
		ServeMux: mux,
		shutdown: shutdown,
		mws:      mws,
	}
}

func (app *App) Handle(pattern string, handler Handler, routemws ...MiddlewareHandler) {
	// Route specific middlewares
	handler = wrapMiddlewares(routemws, handler)
	// Whole app middlewares
	handler = wrapMiddlewares(app.mws, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		if err := handler(r.Context(), w, r); err != nil {
			app.shutdown <- syscall.SIGTERM
		}
	}
	app.ServeMux.HandleFunc(pattern, h)
}
