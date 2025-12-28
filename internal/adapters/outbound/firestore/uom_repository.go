package firestore

import (
	"context"
	// "fmt"

	"github.com/jeffjlins/okra/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	//"errors"
	"github.com/jeffjlins/okra/internal/zerrors"
)

const uomCollection = "uoms"

type UomRepository struct {
	client *Client
}

func NewUomRepository(client *Client) *UomRepository {
	return &UomRepository{
		client: client,
	}
}

// Is it either sentinel error or you need a custom error for every thing? Otherwise how do you distinguish later if not a sentinel error?
// You could have a sentinel error at the very bottom that is static and then wrap it
// But if you need to wrap at the start then it doesn't work as sentinel because it would be different every time
// Instead, you could wrap in an error that has an enum string to identify the error. This just moves the information to a field rather than the type system.


// type UomRepositoryError string

// const (
//     ErrUomValidation      UomRepositoryError = "repo_uom_validation_failure"
//     ErrUomWrite   UomRepositoryError = "repo_uom_write_failure"
// 	ErrUomRead   UomRepositoryError = "repo_uom_read_failure"
//     ErrUomUnmarshall  UomRepositoryError = "repo_uom_unmarshall_failure"
// 	ErrUomDelete  UomRepositoryError = "repo_uom_delete_failure"
// )





// type UOMError struct {
// 	Msg string
// 	Wrapped error
// }
// func (e *UOMError) Error() string {
// 	return e.Msg
// }
// func (e *UOMError) Unwrap() error {
// 	return e.Wrapped
// }

// func ValidationFailedError(err error) error {
// 	return fmt.Errorf("validation failed: %w", err)
// }








func (r *UomRepository) Save(ctx context.Context, uom *domain.Uom) *zerrors.Error[domain.UomRepositoryError] {
	if err := uom.Validate(); err != nil {
		//return fmt.Errorf("validation failed: %w", err)
		return zerrors.New(domain.ErrUomValidation).WithError(err)
	}
	_, err := r.client.Collection(uomCollection).Doc(uom.Id).Set(ctx, uom)
	if err != nil {
		//return fmt.Errorf("failed to save uom %s: %w", uom.Id, err)
		return zerrors.New(domain.ErrUomWrite).WithError(err).With("id", uom.Id)
	}
	return nil
}

func (r *UomRepository) GetByID(ctx context.Context, id string) (*domain.Uom, *zerrors.Error[domain.UomRepositoryError]) {
	doc, err := r.client.Collection(uomCollection).Doc(id).Get(ctx)
	if err != nil {
		// Check if document doesn't exist (NotFound error)
		if status.Code(err) == codes.NotFound {
			return nil, nil // Not found
		}
		//return nil, fmt.Errorf("failed to get uom %s: %w", id, err)
		return nil, zerrors.New(domain.ErrUomRead).WithError(err).With("id", id)
	}

	var uom domain.Uom
	if err := doc.DataTo(&uom); err != nil {
		//return nil, fmt.Errorf("failed to unmarshal uom %s: %w", id, err)
		return nil, zerrors.New(domain.ErrUomUnmarshall).WithError(err).With("id", id)
	}

	return &uom, nil
}

func (r *UomRepository) GetAll(ctx context.Context) ([]*domain.Uom, *zerrors.Error[domain.UomRepositoryError]) {
	docs, err := r.client.Collection(uomCollection).Documents(ctx).GetAll()
	if err != nil {
		//return nil, fmt.Errorf("failed to get all uoms: %w", err)
		return nil, zerrors.New(domain.ErrUomRead).WithError(err)
	}

	uoms := make([]*domain.Uom, 0, len(docs))
	for _, doc := range docs {
		var uom domain.Uom
		if err := doc.DataTo(&uom); err != nil {
			//return nil, fmt.Errorf("failed to unmarshal uom %s: %w", doc.Ref.ID, err)
			return nil, zerrors.New(domain.ErrUomUnmarshall).WithError(err).With("id", doc.Ref.ID)
		}
		uoms = append(uoms, &uom)
	}

	return uoms, nil
}

func (r *UomRepository) Delete(ctx context.Context, id string) *zerrors.Error[domain.UomRepositoryError] {
	_, err := r.client.Collection(uomCollection).Doc(id).Delete(ctx)
	if err != nil {
		//return fmt.Errorf("failed to delete uom %s: %w", id, err)
		return zerrors.New(domain.ErrUomDelete).WithError(err).With("id", id)
	}
	return nil
}
