package usecase

import (
	"context"
	"fmt"

	"github.com/jeffjlins/okra/internal/domain"
)

type UomService struct {
	repo domain.UomRepository
}

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

func (s *UomService) GetUomByID(ctx context.Context, id string) (*domain.Uom, error) {
	uom, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get uom: %w", err)
	}
	if uom == nil {
		return nil, fmt.Errorf("uom with id %s not found", id)
	}
	return uom, nil
}

func (s *UomService) GetAllUoms(ctx context.Context) ([]*domain.Uom, error) {
	uoms, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all uoms: %w", err)
	}
	return uoms, nil
}

func (s *UomService) DeleteUom(ctx context.Context, id string) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error checking for existence of uom: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("uom with id %s not found", id)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete uom: %w", err)
	}
	return nil
}

func (s *UomService) UpdateUom(ctx context.Context, id string, base *domain.BaseUom) (*domain.Uom, error) {
	if err := base.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error checking for existence of uom: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("uom with id %s not found", id)
	}

	uom := &domain.Uom{
		BaseUom: *base,
		Id:      id,
	}

	if err := uom.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := s.repo.Save(ctx, uom); err != nil {
		return nil, fmt.Errorf("failed to update uom: %w", err)
	}

	return uom, nil
}
