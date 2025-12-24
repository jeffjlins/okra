package usecase

import (
	"context"
	"fmt"

	"github.com/jeffjlins/okra/internal/domain"
)

// UomService handles business logic for Uom operations
type UomService struct {
	repo domain.UomRepository
}

// NewUomService creates a new UomService
func NewUomService(repo domain.UomRepository) *UomService {
	return &UomService{
		repo: repo,
	}
}

func (s *UomService) CreateUom(ctx context.Context, base *domain.BaseUom) (*domain.Uom, error) {
	if err := base.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	uom, err := domain.Create(base)
	if err != nil {
		return nil, fmt.Errorf("uom creation failed: %w", err)
	}
	existing, err := s.repo.GetByID(ctx, uom.Id)
	if err != nil {
		return nil, fmt.Errorf("error checking for existence of uom with id %s: %w", uom.Id, err)
	}
	if existing != nil {
		return nil, fmt.Errorf("uom with id %s already exists", uom.Id)
	}

	if err := s.repo.Save(ctx, uom); err != nil {
		return nil, fmt.Errorf("failed to save uom: %w", err)
	}

	return uom, nil
}
