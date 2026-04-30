package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/errs"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/handler/reqcontext"
)

func (h *WebAPIHandler) PutNameByID(w http.ResponseWriter, r *http.Request, pr PathParameters, ctx *reqcontext.RequestContext) {
	id, err := strconv.ParseInt(pr["id"], 10, 64)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, http.StatusBadRequest, "body must be valid JSON with a \"name\" field", err)
		return
	}

	n, err := h.Names.UpdateByID(r.Context(), id, body.Name)
	if errors.Is(err, errs.ErrNotFound) {
		jsonError(w, http.StatusNotFound, "not found", err)
		return
	}
	if errors.Is(err, errs.ErrInvalidInput) {
		jsonError(w, http.StatusBadRequest, "name must be a non-empty string of at most 255 characters", err)
		return
	}
	if err != nil {
		ctx.Logger.Error("PutNameByID service error", "err", err)
		jsonError(w, http.StatusInternalServerError, "internal server error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"id": n.ID, "name": n.Value})
}
