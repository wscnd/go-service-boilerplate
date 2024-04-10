package homeapi

import (
	"context"
	"net/http"
	"github.com/wscnd/go-service-boilerplate/libs/web"
)

func homeHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "ok embedded",
	}
	return web.RespondJSON(ctx, w, status, http.StatusOK)
}

