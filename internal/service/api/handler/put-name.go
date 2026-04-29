package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/service/api/reqcontext"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/service/database"
)

func (h *WebAPIHandler) PutNameByID(w http.ResponseWriter, r *http.Request, pr PathParameters, ctx *reqcontext.RequestContext) {
	id, err := strconv.ParseInt(pr["id"], 10, 64)
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		jsonError(w, "body must be JSON with a non-empty \"name\" field", http.StatusBadRequest)
		return
	}

	if err := h.DB.UpdateNameByID(id, body.Name); errors.Is(err, database.ErrNotFound) {
		jsonError(w, "not found", http.StatusNotFound)
		return
	} else if err != nil {
		ctx.Logger.Error("UpdateNameByID db error", "err", err)
		jsonError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"id": id, "name": body.Name})
}
