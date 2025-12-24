package http

import (
	"net/http"

	"github.com/jeffjlins/okra/internal/usecase"
)

func NewRouter(uomService *usecase.UomService) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /uom", createUomHandler(uomService))
	mux.HandleFunc("GET /uom/{id}", getUomByIDHandler(uomService))
	mux.HandleFunc("GET /uom", getAllUomsHandler(uomService))
	mux.HandleFunc("DELETE /uom/{id}", deleteUomHandler(uomService))
	mux.HandleFunc("PUT /uom/{id}", updateUomHandler(uomService))

	return mux
}
