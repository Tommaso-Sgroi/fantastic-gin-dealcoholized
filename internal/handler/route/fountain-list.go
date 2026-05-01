package routes

import (
	"net/http"
	"strconv"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/domain/fountain"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/handler/reqcontext"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/locationutils"
)

// ListFountains handles GET /fountains/
// Optional query params: latitude, longitude, range (default 10 km).
// When latitude and longitude are both provided, results are filtered by range
// and sorted nearest-first. Otherwise the full unfiltered list is returned.
func (h *WebAPIHandler) ListFountains(ctx *reqcontext.RequestContext) {
	q := ctx.Request.URL.Query()

	var lat, lng *float64
	var rangeKm float64 = 10

	hasLat := q.Has("latitude")
	hasLng := q.Has("longitude")

	if hasLat && hasLng {
		parsedLat, err := locationutils.ParseLatitudeToFloat(q.Get("latitude"))
		if err != nil {
			ctx.Logger.Error("ListFountains invalid latitude", "err", err)
			respondWithJSONError(ctx, http.StatusBadRequest, "invalid latitude")
			return
		}
		parsedLng, err := locationutils.ParseLongitudeToFloat(q.Get("longitude"))
		if err != nil {
			ctx.Logger.Error("ListFountains invalid longitude", "err", err)
			respondWithJSONError(ctx, http.StatusBadRequest, "invalid longitude")
			return
		}
		lat = &parsedLat
		lng = &parsedLng

		if q.Has("range") {
			r, err := strconv.ParseFloat(q.Get("range"), 64)
			if err != nil || r < 1 || r > 200 {
				ctx.Logger.Error("ListFountains invalid range", "err", err)
				respondWithJSONError(ctx, http.StatusBadRequest, "range must be a number between 1 and 200")
				return
			}
			rangeKm = r
		}
	} else if hasLat || hasLng {
		// Only one of the two coordinates was provided — reject the request.
		respondWithJSONError(ctx, http.StatusBadRequest, "latitude and longitude must be provided together")
		return
	}

	fountains, err := h.Fountains.List(ctx.Request.Context(), lat, lng, rangeKm)
	if err != nil {
		ctx.Logger.Error("ListFountains service error", "err", err)
		respondWithJSONError(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	// Always return an array (never null).
	if fountains == nil {
		fountains = make([]fountain.Fountain, 0)
	}

	ctx.JSON(http.StatusOK, fountains)
}
