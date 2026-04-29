package handler

import (
	"net/http"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/service/api/reqcontext"
)

// Liveness checks the API server health. If the server cannot serve requests
// (e.g., database unreachable) it replies with 500; otherwise 200.
func (h *WebAPIHandler) Liveness(w http.ResponseWriter, _ *http.Request, _ PathParameters, ctx *reqcontext.RequestContext) {
	if err := h.DB.Ping(); err != nil {
		ctx.Logger.Error("liveness db ping failed", "err", err)
		jsonError(w, "database unreachable", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
