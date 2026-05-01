package dto

import "github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/domain/name"

type CreateNameRequest struct {
	Name string `json:"name"`
}

type NameResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ToNameResponse maps a domain Name into an API response DTO.
func ToNameResponse(n name.Name) NameResponse {
	return NameResponse{
		ID:   n.ID,
		Name: n.Value,
	}
}
