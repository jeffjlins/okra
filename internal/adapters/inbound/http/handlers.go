package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jeffjlins/okra/internal/domain"
	"github.com/jeffjlins/okra/internal/usecase"
)

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// createUomHandler handles POST /uom requests to create a new Uom
func createUomHandler(uomService *usecase.UomService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// Decode JSON payload into domain.BaseUom
		var base domain.BaseUom
		if err := json.NewDecoder(r.Body).Decode(&base); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		// Create Uom using the service (which handles validation and persistence)
		ctx := r.Context()
		uom, err := uomService.CreateUom(ctx, &base)
		if err != nil {
			log.Printf("Error creating Uom: %v", err)

			// Return appropriate status code based on error type
			statusCode := http.StatusInternalServerError
			errorMsg := "Failed to create Uom"

			// Check for validation or duplicate errors
			errStr := err.Error()
			if strings.Contains(errStr, "validation failed") {
				statusCode = http.StatusBadRequest
				errorMsg = errStr
			} else if strings.Contains(errStr, "already exists") {
				statusCode = http.StatusConflict
				errorMsg = errStr
			}

			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(map[string]string{"error": errorMsg})
			return
		}

		// Return success response with the created Uom
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(uom)
	}
}
