package homeapi

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/wscnd/go-service-boilerplate/apis/errs"
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

func handlerWithError(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100) % 2; n == 0 {
		return errs.New(errors.New("TRUSTED ERROR"), http.StatusBadRequest)
	}
	return fmt.Errorf("some error")
}
