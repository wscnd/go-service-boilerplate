package routes

import (
	"encoding/json"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status string
	}{
		Status: "ok",
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&status)
}
