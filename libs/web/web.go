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
}

func (a *App) Handle(pattern string, handler Handler) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if err := handler(r.Context(), w, r); err != nil {
			a.shutdown <- syscall.SIGTERM
		}
	}
	a.ServeMux.HandleFunc(pattern, h)
}

func NewApp(shutdown chan os.Signal) *App {
	mux := http.NewServeMux()
	return &App{
		ServeMux: mux,
		shutdown: shutdown,
	}
}
