package domain

import (
	"context"
	"github.com/jeffjlins/okra/internal/zerrors"
)

type UomRepository interface {
	Save(ctx context.Context, uom *Uom) *zerrors.Error[UomRepositoryError]
	GetByID(ctx context.Context, id string) (*Uom, *zerrors.Error[UomRepositoryError])
	GetAll(ctx context.Context) ([]*Uom, *zerrors.Error[UomRepositoryError])
	Delete(ctx context.Context, id string) *zerrors.Error[UomRepositoryError]
}

type UomRepositoryError string

const (
    ErrUomValidation      UomRepositoryError = "repo_uom_validation_failure"
    ErrUomWrite   UomRepositoryError = "repo_uom_write_failure"
	ErrUomRead   UomRepositoryError = "repo_uom_read_failure"
    ErrUomUnmarshall  UomRepositoryError = "repo_uom_unmarshall_failure"
	ErrUomDelete  UomRepositoryError = "repo_uom_delete_failure"
)