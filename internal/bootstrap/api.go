package bootstrap

import (
	"net/http"
	"time"

	httpadapter "github.com/jeffjlins/okra/internal/adapters/inbound/http"
)

func NewHTTPServer() *http.Server {
	mux := httpadapter.NewRouter()

	return &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}