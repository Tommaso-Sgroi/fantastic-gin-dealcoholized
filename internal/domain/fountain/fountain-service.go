package fountain

import (
	"context"
	"fmt"
	"sort"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/locationutils"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/repository"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/repository/record"
)

// fountainService is the concrete implementation of ServiceFountain.
type fountainService struct {
	repo repository.Repository
	rtrx repository.Transactor
}

// NewFountainService returns a ServiceFountain backed by the provided Repository.
func NewFountainService(repo repository.Repository, rtrx repository.Transactor) ServiceFountain {
	return &fountainService{repo: repo, rtrx: rtrx}
}

func (s *fountainService) List(ctx context.Context, lat *float64, lng *float64, rangeKm float64) ([]Fountain, error) {
	var rows []record.FountainRecord
	var err error

	if lat != nil && lng != nil {
		// Two-step geo filter:
		// 1. Retrieve all fountains inside a bounding box (fast, SQL-side).
		// 2. Post-filter to a precise circle using the haversine formula (Go-side).
		center := locationutils.Point2D{
			Latitude:  locationutils.Latitude(*lat),
			Longitude: locationutils.Longitude(*lng),
		}
		bbox := locationutils.Squaring(center, rangeKm)
		all, err2 := s.repo.ListFountainsInBoundingBox(ctx,
			float64(bbox.A.Latitude), float64(bbox.B.Latitude),
			float64(bbox.A.Longitude), float64(bbox.B.Longitude),
		)
		if err2 != nil {
			return nil, fmt.Errorf("list (bbox): %w", err2)
		}

		// Circle filter — discard corners of the bounding box.
		for _, r := range all {
			point := locationutils.Point2D{
				Latitude:  locationutils.Latitude(r.Latitude),
				Longitude: locationutils.Longitude(r.Longitude),
			}
			if locationutils.IsPointWithinRadius(center, point, rangeKm) {
				rows = append(rows, r)
			}
		}

		// Sort by distance from the request point (nearest first).
		sort.Slice(rows, func(i, j int) bool {
			pi := locationutils.Point2D{Latitude: locationutils.Latitude(rows[i].Latitude), Longitude: locationutils.Longitude(rows[i].Longitude)}
			pj := locationutils.Point2D{Latitude: locationutils.Latitude(rows[j].Latitude), Longitude: locationutils.Longitude(rows[j].Longitude)}
			return locationutils.HaversineDistance(center, pi) < locationutils.HaversineDistance(center, pj)
		})
	} else {
		rows, err = s.repo.ListFountains(ctx)
		if err != nil {
			return nil, fmt.Errorf("list: %w", err)
		}
	}

	fountains := make([]Fountain, len(rows))
	for i, r := range rows {
		fountains[i] = fountainFromRow(r)
	}
	return fountains, nil
}

func (s *fountainService) Create(ctx context.Context, f Fountain) (Fountain, error) {
	if err := f.Validate(); err != nil {
		return Fountain{}, fmt.Errorf("create: %w", err)
	}
	id, err := s.repo.CreateFountain(ctx, record.FountainRecord{
		Latitude:  f.Latitude,
		Longitude: f.Longitude,
		Status:    f.Status,
	})
	if err != nil {
		return Fountain{}, fmt.Errorf("create: %w", err)
	}
	f.ID = id
	return f, nil
}

func (s *fountainService) Update(ctx context.Context, id uint64, f Fountain) (Fountain, error) {
	if err := f.Validate(); err != nil {
		return Fountain{}, fmt.Errorf("update id=%d: %w", id, err)
	}
	f.ID = id
	if err := s.repo.UpdateFountain(ctx, record.FountainRecord{
		ID:        id,
		Latitude:  f.Latitude,
		Longitude: f.Longitude,
		Status:    f.Status,
	}); err != nil {
		return Fountain{}, fmt.Errorf("update id=%d: %w", id, err)
	}
	return f, nil
}

func (s *fountainService) Delete(ctx context.Context, id uint64) error {
	if err := s.repo.DeleteFountain(ctx, id); err != nil {
		return fmt.Errorf("delete id=%d: %w", id, err)
	}
	return nil
}
