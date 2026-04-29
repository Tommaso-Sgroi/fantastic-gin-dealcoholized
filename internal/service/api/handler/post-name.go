package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/service/api/reqcontext"
)

func (h *WebAPIHandler) PostNewName(w http.ResponseWriter, r *http.Request, _ PathParameters, ctx *reqcontext.RequestContext) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		jsonError(w, "body must be JSON with a non-empty \"name\" field", http.StatusBadRequest)
		return
	}

	id, err := h.DB.CreateName(body.Name)
	if err != nil {
		ctx.Logger.Error("CreateName db error", "err", err)
		jsonError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/name/%d", id))
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"id": id, "name": body.Name})
}
