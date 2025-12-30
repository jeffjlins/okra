package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jeffjlins/okra/internal/domain"
	"github.com/jeffjlins/okra/internal/usecase"
	"github.com/jeffjlins/okra/internal/zerrors"
)

type CustomHandler struct {
	uomService *usecase.UomService
}

func (h *CustomHandler) createUomHandler(w http.ResponseWriter, r *http.Request) *zerrors.Error[HttpError] { // (uomService *usecase.UomService) http.HandlerFunc {
	if r.Method != http.MethodPost {
		//TODO: Provide default messages for status codes?
		return zerrors.New(ErrHttp).With("status_code", http.StatusMethodNotAllowed)
	}

	w.Header().Set("Content-Type", "application/json")

	var base domain.BaseUom
	if err := json.NewDecoder(r.Body).Decode(&base); err != nil {
		//TODO: Add message "Invalid JSON" so it's clear what kind of bad request it is
		return zerrors.New(ErrHttp).WithError(err).With("status_code", http.StatusBadRequest)
	}

	ctx := r.Context()
	uom, err := h.uomService.CreateUom(ctx, &base)
	if err != nil {
		he := zerrors.New(ErrHttp)
		serviceToHttpError(usecase.Generalize(err), he)
		return he
		//TODO: Conflict - "User Id already exists" (shouldn't take id anyway so should be bad request), test this
		//TODO: Bad Request - show all validation errors when that is happening
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(uom)
	return nil
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
