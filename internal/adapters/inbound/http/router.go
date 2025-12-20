package http

import (
	"net/http"

	"github.com/jeffjlins/okra/internal/adapters/outbound/firestore"
)

func NewRouter(repo *firestore.Repository) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /hello", helloHandler)
	mux.HandleFunc("POST /demo", createDemoHandler(repo))

	return mux
}
