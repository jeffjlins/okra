package http

import (
	"net/http"

	"github.com/jeffjlins/okra/internal/adapters/outbound/firestore"
	"github.com/jeffjlins/okra/internal/usecase"
)

func NewRouter(repo *firestore.Repository, uomService *usecase.UomService) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /uom", createUomHandler(uomService))

	return mux
}
