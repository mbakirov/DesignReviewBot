package handlers

import (
	"encoding/json"
	_ "github.com/lib/pq"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	res := struct {
		Message string `json:"message,omitempty"`
	}{
		Message: "OK",
	}

	json.NewEncoder(w).Encode(res)
}