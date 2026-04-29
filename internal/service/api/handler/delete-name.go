package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/service/api/reqcontext"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/service/database"
)

func (h *WebAPIHandler) DeleteNameByID(w http.ResponseWriter, _ *http.Request, pr PathParameters, ctx *reqcontext.RequestContext) {
	id, err := strconv.ParseInt(pr["id"], 10, 64)
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.DB.DeleteNameByID(id); errors.Is(err, database.ErrNotFound) {
		jsonError(w, "not found", http.StatusNotFound)
		return
	} else if err != nil {
		ctx.Logger.Error("DeleteNameByID db error", "err", err)
		jsonError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 — deleted, no body
}
