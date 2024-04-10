package homeapi

import (
	"context"
	"encoding/json"
	"net/http"
)

func homeHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "ok",
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&status)
	return nil
}

