package http

import (
	"net/http"

	"github.com/jeffjlins/okra/internal/zerrors"
)

type HealthAdapter struct {}

func (adp *HealthAdapter) health(w http.ResponseWriter, _ *http.Request) *zerrors.Error[HttpError] {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
	return nil
}