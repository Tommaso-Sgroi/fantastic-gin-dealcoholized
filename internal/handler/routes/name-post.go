package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/errs"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/handler/reqcontext"
)

func (h *WebAPIHandler) PostNewName(w http.ResponseWriter, r *http.Request, _ PathParameters, ctx *reqcontext.RequestContext) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, http.StatusBadRequest, "body must be valid JSON with a \"name\" field")
		return
	}

	n, err := h.Names.Create(r.Context(), body.Name)
	if errors.Is(err, errs.ErrInvalidInput) {
		jsonError(w, http.StatusBadRequest, "name must be a non-empty string of at most 255 characters")
		return
	}
	if err != nil {
		ctx.Logger.Error("PostNewName service error", "err", err)
		jsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/name/%d", n.ID))
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"id": n.ID, "name": n.Value})
}
