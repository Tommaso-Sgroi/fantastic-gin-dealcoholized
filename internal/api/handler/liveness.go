package handler

import (
	"net/http"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/api/reqcontext"
)

// Liveness checks the API server health. If the server cannot serve requests
// (e.g., database unreachable) it replies with 500; otherwise 200.
func (h *Handler) Liveness(w http.ResponseWriter, _ *http.Request, _ PathParameters, ctx *reqcontext.RequestContext) {
	if err := h.DB.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
