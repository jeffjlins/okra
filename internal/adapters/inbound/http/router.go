package http

import (
	"net/http"
	"os"

	"github.com/jeffjlins/okra/internal/usecase"
	"github.com/jeffjlins/okra/internal/zerrors"
)

type APIFunc func(w http.ResponseWriter, r *http.Request) *zerrors.Error[HttpError]

type Router struct {
	logFormat         string
	includeStackTrace bool
}

func (router *Router) make(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			if router.logFormat == "json" {
				err.SlogErrorJSON(os.Stderr)
			} else {
				err.LogErrorText(os.Stderr)
			}
			err.WriteErrorHttpJSON(w, router.includeStackTrace)
			return
		}
	}
}

func NewRouter(uomService *usecase.UomService, logFormat string, includeStackTrace bool) *http.ServeMux {
	router := &Router{
		logFormat:         logFormat,
		includeStackTrace: includeStackTrace,
	}
	mux := http.NewServeMux()

	healthAdp := HealthAdapter{}
	mux.HandleFunc("GET /health", router.make(healthAdp.health))
	uomAdp := UomAdapter{uomService}
	mux.HandleFunc("POST /uom", router.make(uomAdp.create))
	mux.HandleFunc("GET /uom/{id}", router.make(uomAdp.getOneById))
	mux.HandleFunc("GET /uom", router.make(uomAdp.getAll))
	mux.HandleFunc("DELETE /uom/{id}", router.make(uomAdp.delete))
	mux.HandleFunc("PUT /uom/{id}", router.make(uomAdp.update))

	return mux
}
