package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/errs"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/handler/reqcontext"
)

// DeleteFountain handles DELETE /fountains/:id
func (h *WebAPIHandler) DeleteFountain(ctx *reqcontext.RequestContext) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Logger.Error("DeleteFountain invalid id", "err", err)
		respondWithJSONError(ctx, http.StatusBadRequest, "id must be a positive integer")
		return
	}

	err = h.Fountains.Delete(ctx.Request.Context(), id)
	if errors.Is(err, errs.ErrNotFound) {
		ctx.Logger.Error("DeleteFountain not found", "id", id)
		respondWithJSONError(ctx, http.StatusNotFound, "fountain not found")
		return
	}
	if err != nil {
		ctx.Logger.Error("DeleteFountain service error", "id", id, "err", err)
		respondWithJSONError(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	ctx.Status(http.StatusNoContent)
}
