package usecase

import (
	"github.com/jeffjlins/okra/internal/zerrors"
)

type ServiceError string

const (
	ErrSrvNotFound ServiceError = "service.not_found"
	ErrSrvInput    ServiceError = "service.bad_input"
	ErrSrvSystem   ServiceError = "service.system_failure"
	ErrSrvConflict ServiceError = "service.conflict"
)

type AnyServiceError interface {
	UomServiceError // | ConvServiceError | FoodServiceError
}

func GeneralizeError[T AnyServiceError](err *zerrors.Error[T]) *zerrors.Error[ServiceError] {
	switch any(err.Code()).(type) {
	case UomServiceError:
		if err2, ok := any(err).(*zerrors.Error[UomServiceError]); ok {
			return generalizeErrorUom(err2)
		}
	}
	return zerrors.New(ErrSrvSystem)
}
