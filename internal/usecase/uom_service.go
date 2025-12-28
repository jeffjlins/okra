package usecase

import (
	"context"
	// "fmt"

	"github.com/jeffjlins/okra/internal/domain"
	"github.com/jeffjlins/okra/internal/zerrors"
)

type UomService struct {
	repo domain.UomRepository
}

func NewUomService(repo domain.UomRepository) *UomService {
	return &UomService{
		repo: repo,
	}
}

type UomServiceError string

const (
    ErrUomValidation      UomServiceError = "serv_uom_validation_failure"
    ErrUomDbCreate   UomServiceError = "serv_uom_db_create_failure"
	ErrUomDbRead   UomServiceError = "serv_uom_db_read_failure"
	ErrUomDbWrite   UomServiceError = "serv_uom_db_write_failure"
	ErrUomDbDelete   UomServiceError = "serv_uom_db_delete_failure"
	ErrUomDbExistenceCheck   UomServiceError = "serv_uom_db_existence_check_failure"
	ErrUomIdOverwrite   UomServiceError = "serv_uom_id_overwrite_not_allowed"
	ErrUomNotFound   UomServiceError = "serv_uom_not_found"
)

func generalizeUom(err *zerrors.Error[UomServiceError]) *zerrors.Error[ServiceError] {
	switch err.Code() {
	case ErrUomDbCreate, ErrUomDbRead, ErrUomDbWrite, ErrUomDbDelete, ErrUomDbExistenceCheck:
		return zerrors.New(ErrSrvSystem).WithError(err)
	case ErrUomNotFound:
		return zerrors.New(ErrSrvNotFound).WithError(err)
	case ErrUomIdOverwrite:
		return zerrors.New(ErrSrvConflict).WithError(err)
	case ErrUomValidation:
		return zerrors.New(ErrSrvInput).WithError(err)
	}
	return zerrors.New(ErrSrvSystem) // unknown error
}



func (s *UomService) CreateUom(ctx context.Context, base *domain.BaseUom) (*domain.Uom, *zerrors.Error[UomServiceError]) {
	if err := base.Validate(); err != nil {
		// return nil, fmt.Errorf("validation failed: %w", err)
		return nil, zerrors.New(ErrUomValidation).WithError(err)
	}
	uom, err := domain.Create(base)
	if err != nil {
		// return nil, fmt.Errorf("uom creation failed: %w", err)
		return nil, zerrors.New(ErrUomDbCreate).WithError(err)
	}
	existing, err := s.repo.GetByID(ctx, uom.Id)
	if err != nil {
		//return nil, fmt.Errorf("error checking for existence of uom with id %s: %w", uom.Id, err)
		return nil, zerrors.New(ErrUomDbExistenceCheck).WithError(err).With("id", uom.Id)
	}
	if existing != nil {
		// return nil, fmt.Errorf("uom with id %s already exists", uom.Id)
		return nil, zerrors.New(ErrUomDbExistenceCheck).With("id", uom.Id)
	}

	if err := s.repo.Save(ctx, uom); err != nil {
		// return nil, fmt.Errorf("failed to save uom: %w", err)
		return nil, zerrors.New(ErrUomDbWrite).WithError(err)
	}

	return uom, nil
}

func (s *UomService) GetUomByID(ctx context.Context, id string) (*domain.Uom, *zerrors.Error[UomServiceError]) {
	uom, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// return nil, fmt.Errorf("failed to get uom: %w", err)
		return nil, zerrors.New(ErrUomDbRead).WithError(err)
	}
	if uom == nil {
		// return nil, fmt.Errorf("uom with id %s not found", id)
		return nil, zerrors.New(ErrUomNotFound)
	}
	return uom, nil
}

func (s *UomService) GetAllUoms(ctx context.Context) ([]*domain.Uom, *zerrors.Error[UomServiceError]) {
	uoms, err := s.repo.GetAll(ctx)
	if err != nil {
		// return nil, fmt.Errorf("failed to get all uoms: %w", err)
		zerrors.New(ErrUomDbRead).WithError(err)
	}
	return uoms, nil
}

func (s *UomService) DeleteUom(ctx context.Context, id string) *zerrors.Error[UomServiceError] {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// return fmt.Errorf("error checking for existence of uom: %w", err)
		return zerrors.New(ErrUomDbExistenceCheck).WithError(err)
	}
	if existing == nil {
		// return fmt.Errorf("uom with id %s not found", id)
		return zerrors.New(ErrUomNotFound)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		// return fmt.Errorf("failed to delete uom: %w", err)
		return zerrors.New(ErrUomDbDelete).WithError(err)
	}
	return nil
}

func (s *UomService) UpdateUom(ctx context.Context, id string, base *domain.BaseUom) (*domain.Uom, *zerrors.Error[UomServiceError]) {
	if err := base.Validate(); err != nil {
		// return nil, fmt.Errorf("validation failed: %w", err)
		return nil, zerrors.New(ErrUomValidation).WithError(err)
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// return nil, fmt.Errorf("error checking for existence of uom: %w", err)
		return nil, zerrors.New(ErrUomDbExistenceCheck).WithError(err)
	}
	if existing == nil {
		// return nil, fmt.Errorf("uom with id %s not found", id)
		return nil, zerrors.New(ErrUomNotFound)
	}

	uom := &domain.Uom{
		BaseUom: *base,
		Id:      id,
	}

	if err := uom.Validate(); err != nil {
		// return nil, fmt.Errorf("validation failed: %w", err)
		return nil, zerrors.New(ErrUomValidation).WithError(err)
	}

	if err := s.repo.Save(ctx, uom); err != nil {
		// return nil, fmt.Errorf("failed to update uom: %w", err)
		return nil, zerrors.New(ErrUomDbWrite).WithError(err)
	}

	return uom, nil
}
