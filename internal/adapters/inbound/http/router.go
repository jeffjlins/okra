package http

import (
	"net/http"
	"os"

	"github.com/jeffjlins/okra/internal/usecase"
	"github.com/jeffjlins/okra/internal/zerrors"
)

type APIFunc func(w http.ResponseWriter, r *http.Request) *zerrors.Error[HttpError]

type Router struct {
	logFormat string
}

func make(router *Router, f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			if router.logFormat == "json" {
				err.SlogErrorJSON(os.Stderr)
			} else {
				err.LogErrorText(os.Stderr)
			}
			err.WriteErrorHttpJSON(w)
			return
		}
	}
}

func NewRouter(uomService *usecase.UomService, logFormat string) *http.ServeMux {
	router := &Router{logFormat: logFormat}
	mux := http.NewServeMux()

	var ch CustomHandler = CustomHandler{uomService}
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /uom", make(router, ch.createUomHandler))
	mux.HandleFunc("GET /uom/{id}", getUomByIDHandler(uomService))
	mux.HandleFunc("GET /uom", getAllUomsHandler(uomService))
	mux.HandleFunc("DELETE /uom/{id}", deleteUomHandler(uomService))
	mux.HandleFunc("PUT /uom/{id}", updateUomHandler(uomService))

	return mux
}
