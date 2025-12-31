package http

import (
	"encoding/json"
	"net/http"

	"github.com/jeffjlins/okra/internal/domain"
	"github.com/jeffjlins/okra/internal/usecase"
	"github.com/jeffjlins/okra/internal/zerrors"
)

type UomAdapter struct {
	uomService *usecase.UomService
}

func (adp *UomAdapter) create(w http.ResponseWriter, r *http.Request) *zerrors.Error[HttpError] {
	if r.Method != http.MethodPost {
		//TODO: Provide default messages for status codes?
		return zerrors.New(ErrHttp).With("status_code", http.StatusMethodNotAllowed)
	}

	w.Header().Set("Content-Type", "application/json")

	var base domain.BaseUom
	if err := json.NewDecoder(r.Body).Decode(&base); err != nil {
		return zerrors.New(ErrHttp).WithError(err).With("status_code", http.StatusBadRequest).With("reason", "Invalid JSON")
	}

	ctx := r.Context()
	uom, err := adp.uomService.CreateUom(ctx, &base)
	if err != nil {
		e := zerrors.New(ErrHttp)
		serviceToHttpError(usecase.GeneralizeError(err), e)
		return e
		//TODO: Conflict - "User Id already exists" (shouldn't take id anyway so should be bad request), test this
		//TODO: Bad Request - show all validation errors when that is happening
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(uom)
	return nil
}

func (adp *UomAdapter) getOneById(w http.ResponseWriter, r *http.Request) *zerrors.Error[HttpError] {
	if r.Method != http.MethodGet {
		return zerrors.New(ErrHttp).With("status_code", http.StatusMethodNotAllowed)
	}

	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		return zerrors.New(ErrHttp).With("status_code", http.StatusBadRequest).With("reason", "id is required")
	}

	ctx := r.Context()
	uom, err := adp.uomService.GetUomByID(ctx, id)
	if err != nil {
		e := zerrors.New(ErrHttp)
		serviceToHttpError(usecase.GeneralizeError(err), e)
		return e
	}

	json.NewEncoder(w).Encode(uom)
	return nil
}

func (adp *UomAdapter) getAll(w http.ResponseWriter, r *http.Request) *zerrors.Error[HttpError] {
	if r.Method != http.MethodGet {
		return zerrors.New(ErrHttp).With("status_code", http.StatusMethodNotAllowed)
	}

	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	uoms, err := adp.uomService.GetAllUoms(ctx)
	if err != nil {
		e := zerrors.New(ErrHttp)
		serviceToHttpError(usecase.GeneralizeError(err), e)
		return e
	}

	json.NewEncoder(w).Encode(uoms)
	return nil
}

func (adp *UomAdapter) delete(w http.ResponseWriter, r *http.Request) *zerrors.Error[HttpError] {
	if r.Method != http.MethodDelete {
		return zerrors.New(ErrHttp).With("status_code", http.StatusMethodNotAllowed)
	}

	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return zerrors.New(ErrHttp).With("status_code", http.StatusBadRequest).With("reason", "id is required")
	}

	ctx := r.Context()
	if err := adp.uomService.DeleteUom(ctx, id); err != nil {
		e := zerrors.New(ErrHttp)
		serviceToHttpError(usecase.GeneralizeError(err), e)
		return e
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (adp *UomAdapter) update(w http.ResponseWriter, r *http.Request) *zerrors.Error[HttpError] {
	if r.Method != http.MethodPut {
		return zerrors.New(ErrHttp).With("status_code", http.StatusMethodNotAllowed)
	}

	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		return zerrors.New(ErrHttp).With("status_code", http.StatusBadRequest).With("reason", "id is required")
	}

	var base domain.BaseUom
	if err := json.NewDecoder(r.Body).Decode(&base); err != nil {
		return zerrors.New(ErrHttp).With("status_code", http.StatusBadRequest).With("reason", "invalid json")
	}

	ctx := r.Context()
	uom, err := adp.uomService.UpdateUom(ctx, id, &base)
	if err != nil {
		e := zerrors.New(ErrHttp)
		serviceToHttpError(usecase.GeneralizeError(err), e)
		return e
	}

	json.NewEncoder(w).Encode(uom)
	return nil
}
