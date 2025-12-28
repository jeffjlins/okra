package http

import (
	// "encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jeffjlins/okra/internal/usecase"
	"github.com/jeffjlins/okra/internal/zerrors"
)

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func Make(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			var e *zerrors.Error[HttpError]
			if !errors.As(err, &e) {
				e = zerrors.New(ErrHttp).With("status_code", http.StatusInternalServerError)
			}
			writeErrorJSON(w, e)
			return
		}
	}
}

func writeErrorJSON(w http.ResponseWriter, e *zerrors.Error[HttpError]) {
	var code int
	if c, ok := e.Get("status_code"); ok {
		code = c.(int)
	} else {
		code = http.StatusInternalServerError
	}
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	logger := slog.New(slog.NewJSONHandler(w, nil))
	logger.Error("Failed to find user", slog.Any("error", e))
}

func NewRouter(uomService *usecase.UomService) *http.ServeMux {
	mux := http.NewServeMux()

	var ch CustomHandler = CustomHandler{uomService}
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /uom", Make(ch.createUomHandler))
	mux.HandleFunc("GET /uom/{id}", getUomByIDHandler(uomService))
	mux.HandleFunc("GET /uom", getAllUomsHandler(uomService))
	mux.HandleFunc("DELETE /uom/{id}", deleteUomHandler(uomService))
	mux.HandleFunc("PUT /uom/{id}", updateUomHandler(uomService))

	return mux
}
