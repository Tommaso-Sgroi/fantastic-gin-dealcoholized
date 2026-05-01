package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/domain/fountain"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/errs"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/handler/reqcontext"
)

// CreateFountain handles POST /fountains/
func (h *WebAPIHandler) CreateFountain(ctx *reqcontext.RequestContext) {
	var f fountain.Fountain
	if err := json.NewDecoder(ctx.Request.Body).Decode(&f); err != nil {
		ctx.Logger.Error("CreateFountain invalid request body", "err", err)
		respondWithJSONError(ctx, http.StatusBadRequest, "body must be valid JSON")
		return
	}

	created, err := h.Fountains.Create(ctx.Request.Context(), f)
	if errors.Is(err, errs.ErrInvalidInput) {
		ctx.Logger.Error("CreateFountain invalid input", "err", err)
		respondWithJSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		ctx.Logger.Error("CreateFountain service error", "err", err)
		respondWithJSONError(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	ctx.JSON(http.StatusCreated, created)
}
