package handler

import (
	"encoding/json"
	"net/http"
)

// jsonError writes a consistent JSON error response: {"error": "<msg>"}
func jsonError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
