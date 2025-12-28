package usecase

import (
	"github.com/jeffjlins/okra/internal/zerrors"
)

type ServiceError string

const (
	ErrSrvNotFound ServiceError = "not_found"
	ErrSrvInput    ServiceError = "bad_input"
	ErrSrvSystem   ServiceError = "system_failure"
	ErrSrvConflict ServiceError = "conflict"
)

type AnyServiceError interface {
	UomServiceError // | ConvServiceError | FoodServiceError
}

func Generalize[T AnyServiceError](err *zerrors.Error[T]) *zerrors.Error[ServiceError] {
	switch any(err.Code()).(type) {
	case UomServiceError:
		if err2, ok := any(err).(*zerrors.Error[UomServiceError]); ok {
			return generalizeUom(err2)
		}
	}
	return zerrors.New(ErrSrvSystem)
}
