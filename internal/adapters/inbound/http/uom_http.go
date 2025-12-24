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

func createUomHandler(uomService *usecase.UomService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		var base domain.BaseUom
		if err := json.NewDecoder(r.Body).Decode(&base); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		uom, err := uomService.CreateUom(ctx, &base)
		if err != nil {
			log.Printf("Error creating Uom: %v", err)

			statusCode := http.StatusInternalServerError
			errorMsg := "Failed to create Uom"

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

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(uom)
	}
}

func getUomByIDHandler(uomService *usecase.UomService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		uom, err := uomService.GetUomByID(ctx, id)
		if err != nil {
			log.Printf("Error getting Uom: %v", err)

			statusCode := http.StatusInternalServerError
			if strings.Contains(err.Error(), "not found") {
				statusCode = http.StatusNotFound
			}

			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(uom)
	}
}

func getAllUomsHandler(uomService *usecase.UomService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		ctx := r.Context()
		uoms, err := uomService.GetAllUoms(ctx)
		if err != nil {
			log.Printf("Error getting all Uoms: %v", err)
			http.Error(w, "Failed to get Uoms", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(uoms)
	}
}

func deleteUomHandler(uomService *usecase.UomService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		if err := uomService.DeleteUom(ctx, id); err != nil {
			log.Printf("Error deleting Uom: %v", err)

			statusCode := http.StatusInternalServerError
			if strings.Contains(err.Error(), "not found") {
				statusCode = http.StatusNotFound
			}

			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func updateUomHandler(uomService *usecase.UomService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		var base domain.BaseUom
		if err := json.NewDecoder(r.Body).Decode(&base); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		uom, err := uomService.UpdateUom(ctx, id, &base)
		if err != nil {
			log.Printf("Error updating Uom: %v", err)

			statusCode := http.StatusInternalServerError
			errorMsg := "Failed to update Uom"

			errStr := err.Error()
			if strings.Contains(errStr, "not found") {
				statusCode = http.StatusNotFound
				errorMsg = errStr
			} else if strings.Contains(errStr, "validation failed") {
				statusCode = http.StatusBadRequest
				errorMsg = errStr
			}

			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(map[string]string{"error": errorMsg})
			return
		}

		json.NewEncoder(w).Encode(uom)
	}
}
