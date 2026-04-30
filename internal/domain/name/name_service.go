package name

import (
	"context"
	"fmt"
	"strings"
)

// NameRepository is the persistence contract that the domain requires.
// It is implemented by the database layer. The domain never imports database/sql
// directly — the repository hides all storage details.
type NameRepository interface {
	GetNameByID(id int64) (string, error)
	CreateName(name string) (int64, error)
	UpdateNameByID(id int64, name string) error
	DeleteNameByID(id int64) error
}

// NameService defines the business operations available on Name entities.
// All methods accept a context.Context for cancellation and tracing propagation.
type NameService interface {
	// GetByID returns the Name with the given id.
	// Returns ErrNotFound when no such entity exists.
	GetByID(ctx context.Context, id int64) (Name, error)

	// Create validates and persists a new Name, returning the created entity.
	// Returns ErrInvalidInput when the name is empty or exceeds the allowed length.
	Create(ctx context.Context, value string) (Name, error)

	// UpdateByID replaces the value of the Name identified by id.
	// Returns ErrNotFound when no such entity exists.
	// Returns ErrInvalidInput when the new value fails validation.
	UpdateByID(ctx context.Context, id int64, value string) (Name, error)

	// DeleteByID removes the Name identified by id.
	// Returns ErrNotFound when no such entity exists.
	DeleteByID(ctx context.Context, id int64) error
}

const maxNameLength = 255

// nameService is the concrete implementation of NameService.
type nameService struct {
	repo NameRepository
}

// NewNameService returns a NameService backed by the provided NameRepository.
func NewNameService(repo NameRepository) NameService {
	return &nameService{repo: repo}
}

func (s *nameService) GetByID(_ context.Context, id int64) (Name, error) {
	value, err := s.repo.GetNameByID(id)
	if err != nil {
		return Name{}, fmt.Errorf("cannot find %d in db: %w", id, err)
	}
	return Name{ID: id, Value: value}, nil
}

func (s *nameService) Create(_ context.Context, value string) (Name, error) {
	if err := validateName(value); err != nil {
		return Name{}, fmt.Errorf("Create: %w", err)
	}
	id, err := s.repo.CreateName(value)
	if err != nil {
		return Name{}, fmt.Errorf("Create: %w", err)
	}
	return Name{ID: id, Value: value}, nil
}

func (s *nameService) UpdateByID(_ context.Context, id int64, value string) (Name, error) {
	if err := validateName(value); err != nil {
		return Name{}, fmt.Errorf("UpdateByID id=%d: %w", id, err)
	}
	if err := s.repo.UpdateNameByID(id, value); err != nil {
		return Name{}, fmt.Errorf("UpdateByID id=%d: %w", id, err)
	}
	return Name{ID: id, Value: value}, nil
}

func (s *nameService) DeleteByID(_ context.Context, id int64) error {
	if err := s.repo.DeleteNameByID(id); err != nil {
		return fmt.Errorf("DeleteByID id=%d: %w", id, err)
	}
	return nil
}

// validateName enforces domain-level invariants on a Name value.
func validateName(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("name cannot be empty or whitespace: %w", ErrInvalidInput)
	}
	if len(value) > maxNameLength {
		return fmt.Errorf("name exceeds %d characters: %w", maxNameLength, ErrInvalidInput)
	}
	return nil
}
