package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/errs"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/handler/reqcontext"
)

func (h *WebAPIHandler) DeleteNameByID(w http.ResponseWriter, r *http.Request, pr PathParameters, ctx *reqcontext.RequestContext) {
	id, err := strconv.ParseInt(pr["id"], 10, 64)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id", err)
		return
	}

	if err := h.Names.DeleteByID(r.Context(), id); errors.Is(err, errs.ErrNotFound) {
		jsonError(w, http.StatusNotFound, "not found", err)
		return
	} else if err != nil {
		ctx.Logger.Error("DeleteNameByID service error", "err", err)
		jsonError(w, http.StatusInternalServerError, "internal server error", err)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 — deleted, no-body
}
