package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/service/api/reqcontext"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/service/database"
)

func (h *WebAPIHandler) GetNameByID(w http.ResponseWriter, hr *http.Request, pr PathParameters, ctx *reqcontext.RequestContext) {
	id, err := strconv.ParseInt(pr["id"], 10, 64)
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	name, err := h.DB.GetNameByID(id)
	if errors.Is(err, database.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
		jsonError(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		ctx.Logger.Error("GetNameByID db error", "err", err)
		jsonError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"id": id, "name": name})
}
