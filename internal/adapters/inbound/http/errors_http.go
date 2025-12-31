package http

import (
	"net/http"

	"github.com/jeffjlins/okra/internal/usecase"
	"github.com/jeffjlins/okra/internal/zerrors"
)

type HttpError string

const ErrHttp HttpError = "adapter_http.error"

func serviceToHttpError(se *zerrors.Error[usecase.ServiceError], he *zerrors.Error[HttpError]) {
	switch se.Code() {
	case usecase.ErrSrvConflict:
		he.With("status_code", http.StatusConflict).WithError(se)
	case usecase.ErrSrvInput:
		he.With("status_code", http.StatusBadRequest).WithError(se)
	case usecase.ErrSrvNotFound:
		he.With("status_code", http.StatusNotFound).WithError(se)
	case usecase.ErrSrvSystem:
		he.With("status_code", http.StatusInternalServerError).WithError(se)
	}
}
