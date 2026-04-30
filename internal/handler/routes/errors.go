package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// jsonError writes a consistent JSON error response: {"error": "<msg>"}
func jsonError(w http.ResponseWriter, status int, msg string, err error) {
	if err != nil {
		msg = fmt.Sprintf("%s: %s", msg, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
