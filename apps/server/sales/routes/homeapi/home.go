package homeapi

import (
	"encoding/json"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status string
	}{
		Status: "ok",
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&status)
}
