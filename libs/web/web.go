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
			// TODO: handle errors
			app.SignalShutdown()
		}
	}
	app.ServeMux.HandleFunc(pattern, h)
}

// SignalShutdown is used to gracefully shutdown the app when integrity issue
// is identified. It means that the error went through the error handler and
// is propagating.
func (app *App) SignalShutdown() {
	app.shutdown <- syscall.SIGTERM
}
