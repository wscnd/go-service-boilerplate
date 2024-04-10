package web

import (
	"context"
	"net/http"
	"os"
)

type App struct {
	*http.ServeMux
	shutdown chan os.Signal
}

func NewApp(shutdown chan os.Signal) *App {
	mux := http.NewServeMux()
	return &App{
		ServeMux: mux,
		shutdown: shutdown,
	}
}
