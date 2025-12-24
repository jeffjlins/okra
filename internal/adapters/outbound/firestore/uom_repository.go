package firestore

import (
	"context"
	"fmt"

	"github.com/jeffjlins/okra/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (r *UomRepository) Save(ctx context.Context, uom *domain.Uom) error {
	if err := uom.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	_, err := r.client.Collection(uomCollection).Doc(uom.Id).Set(ctx, uom)
	if err != nil {
		return fmt.Errorf("failed to save uom %s: %w", uom.Id, err)
	}
	return nil
}

func (r *UomRepository) GetByID(ctx context.Context, id string) (*domain.Uom, error) {
	doc, err := r.client.Collection(uomCollection).Doc(id).Get(ctx)
	if err != nil {
		// Check if document doesn't exist (NotFound error)
		if status.Code(err) == codes.NotFound {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to get uom %s: %w", id, err)
	}

	var uom domain.Uom
	if err := doc.DataTo(&uom); err != nil {
		return nil, fmt.Errorf("failed to unmarshal uom %s: %w", id, err)
	}

	return &uom, nil
}
