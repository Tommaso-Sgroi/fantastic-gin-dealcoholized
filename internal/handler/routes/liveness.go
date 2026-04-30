package routes

import (
	"net/http"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/handler/reqcontext"
)

// Liveness checks the application health. Replies 500 if the database is
// unreachable, 200 otherwise.
func (h *WebAPIHandler) Liveness(w http.ResponseWriter, r *http.Request, _ PathParameters, ctx *reqcontext.RequestContext) {
	if err := h.Health.Liveness(r.Context()); err != nil {
		ctx.Logger.Error("liveness check failed", "err", err)
		jsonError(w, http.StatusInternalServerError, "database unreachable", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
