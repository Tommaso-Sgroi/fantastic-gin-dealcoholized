package dto

import "github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/domain/fountain"

// FountainRequest is the JSON body accepted by POST /fountains/ and PUT /fountains/:id.
// The ID is never provided by the client — it comes from the DB (create) or the URL path (update).
type FountainRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Status    string  `json:"status"`
}

// ToDomain converts a FountainRequest DTO into a domain Fountain.
// ID is left zero — callers supply it from the URL path when needed (e.g. update).
func (r FountainRequest) ToDomain() fountain.Fountain {
	return fountain.Fountain{
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
		Status:    r.Status,
	}
}

// FountainResponse is the JSON body returned by all single-fountain endpoints.
// Every field is always populated so clients never see partial objects.
// Use NewFountainResponse to construct it — never build the struct literal directly.
type FountainResponse struct {
	ID        uint64  `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Status    string  `json:"status"`
}

// NewFountainResponse constructs a FountainResponse from a domain Fountain.
func NewFountainResponse(f fountain.Fountain) FountainResponse {
	return FountainResponse{
		ID:        f.ID,
		Latitude:  f.Latitude,
		Longitude: f.Longitude,
		Status:    f.Status,
	}
}

// FountainListResponse is the JSON body returned by GET /fountains/.
// Use NewFountainListResponse to construct it — never build the struct literal directly.
type FountainListResponse struct {
	Fountains []FountainResponse `json:"fountains"`
}

// NewFountainListResponse constructs a FountainListResponse from a slice of domain Fountains.
// The Fountains field is always a non-nil slice so clients never receive a JSON null.
func NewFountainListResponse(fs []fountain.Fountain) FountainListResponse {
	items := make([]FountainResponse, len(fs))
	for i, f := range fs {
		items[i] = NewFountainResponse(f)
	}
	return FountainListResponse{Fountains: items}
}
