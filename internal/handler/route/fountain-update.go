package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/domain/fountain"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/errs"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/handler/reqcontext"
)

// UpdateFountain handles PUT /fountains/:id
func (h *WebAPIHandler) UpdateFountain(ctx *reqcontext.RequestContext) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Logger.Error("UpdateFountain invalid id", "err", err)
		respondWithJSONError(ctx, http.StatusBadRequest, "id must be a positive integer")
		return
	}

	var f fountain.Fountain
	if err = json.NewDecoder(ctx.Request.Body).Decode(&f); err != nil {
		ctx.Logger.Error("UpdateFountain invalid request body", "err", err)
		respondWithJSONError(ctx, http.StatusBadRequest, "body must be valid JSON")
		return
	}

	updated, err := h.Fountains.Update(ctx.Request.Context(), id, f)
	if errors.Is(err, errs.ErrNotFound) {
		ctx.Logger.Error("UpdateFountain not found", "id", id)
		respondWithJSONError(ctx, http.StatusNotFound, "fountain not found")
		return
	}
	if errors.Is(err, errs.ErrInvalidInput) {
		ctx.Logger.Error("UpdateFountain invalid input", "err", err)
		respondWithJSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		ctx.Logger.Error("UpdateFountain service error", "id", id, "err", err)
		respondWithJSONError(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, updated)
}
