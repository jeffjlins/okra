package usecase

import (
	"context"

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
	ErrUomValidation       UomServiceError = "serv_uom_validation_failure"
	ErrUomDbCreate         UomServiceError = "serv_uom_db_create_failure"
	ErrUomDbRead           UomServiceError = "serv_uom_db_read_failure"
	ErrUomDbWrite          UomServiceError = "serv_uom_db_write_failure"
	ErrUomDbDelete         UomServiceError = "serv_uom_db_delete_failure"
	ErrUomDbExistenceCheck UomServiceError = "serv_uom_db_existence_check_failure"
	ErrUomIdOverwrite      UomServiceError = "serv_uom_id_overwrite_not_allowed"
	ErrUomNotFound         UomServiceError = "serv_uom_not_found"
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
		return nil, zerrors.New(ErrUomValidation).WithError(err)
	}
	uom, err := domain.Create(base)
	if err != nil {
		return nil, zerrors.New(ErrUomDbCreate).WithError(err)
	}
	existing, err2 := s.repo.GetByID(ctx, uom.Id)
	if err2 != nil {
		return nil, zerrors.New(ErrUomDbExistenceCheck).WithError(err2).With("id", uom.Id)
	}
	if existing != nil {
		return nil, zerrors.New(ErrUomDbExistenceCheck).With("id", uom.Id)
	}

	if err2 := s.repo.Save(ctx, uom); err2 != nil {
		return nil, zerrors.New(ErrUomDbWrite).WithError(err2).With("id", uom.Id)
	}

	return uom, nil
}

func (s *UomService) GetUomByID(ctx context.Context, id string) (*domain.Uom, *zerrors.Error[UomServiceError]) {
	uom, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, zerrors.New(ErrUomDbRead).WithError(err).With("id", id)
	}
	if uom == nil {
		return nil, zerrors.New(ErrUomNotFound).With("id", id)
	}
	return uom, nil
}

func (s *UomService) GetAllUoms(ctx context.Context) ([]*domain.Uom, *zerrors.Error[UomServiceError]) {
	uoms, err := s.repo.GetAll(ctx)
	if err != nil {
		zerrors.New(ErrUomDbRead).WithError(err)
	}
	return uoms, nil
}

func (s *UomService) DeleteUom(ctx context.Context, id string) *zerrors.Error[UomServiceError] {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return zerrors.New(ErrUomDbExistenceCheck).WithError(err).With("id", id)
	}
	if existing == nil {
		return zerrors.New(ErrUomNotFound).With("id", id)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return zerrors.New(ErrUomDbDelete).WithError(err).With("id", id)
	}
	return nil
}

func (s *UomService) UpdateUom(ctx context.Context, id string, base *domain.BaseUom) (*domain.Uom, *zerrors.Error[UomServiceError]) {
	if err := base.Validate(); err != nil {
		return nil, zerrors.New(ErrUomValidation).WithError(err).With("id", id)
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, zerrors.New(ErrUomDbExistenceCheck).WithError(err).With("id", id)
	}
	if existing == nil {
		return nil, zerrors.New(ErrUomNotFound).With("id", id)
	}

	uom := &domain.Uom{
		BaseUom: *base,
		Id:      id,
	}

	if err := uom.Validate(); err != nil {
		return nil, zerrors.New(ErrUomValidation).WithError(err).With("id", id)
	}

	if err := s.repo.Save(ctx, uom); err != nil {
		return nil, zerrors.New(ErrUomDbWrite).WithError(err).With("id", id)
	}

	return uom, nil
}
