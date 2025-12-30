package firestore

import (
	"context"
	
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jeffjlins/okra/internal/domain"
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

func (r *UomRepository) Save(ctx context.Context, uom *domain.Uom) *zerrors.Error[domain.UomRepositoryError] {
	if err := uom.Validate(); err != nil {
		return zerrors.New(domain.ErrUomValidation).WithError(err)
	}
	_, err := r.client.Collection(uomCollection).Doc(uom.Id).Set(ctx, uom)
	if err != nil {
		return zerrors.New(domain.ErrUomWrite).WithError(err).With("id", uom.Id)
	}
	return nil
}

func (r *UomRepository) GetByID(ctx context.Context, id string) (*domain.Uom, *zerrors.Error[domain.UomRepositoryError]) {
	doc, err := r.client.Collection(uomCollection).Doc(id).Get(ctx)
	if err != nil {
		// Check if document doesn't exist (NotFound error)
		if status.Code(err) == codes.NotFound {
			return nil, nil //TODO: Not found Error
		}
		return nil, zerrors.New(domain.ErrUomRead).WithError(err).With("id", id)
	}

	var uom domain.Uom
	if err := doc.DataTo(&uom); err != nil {
		return nil, zerrors.New(domain.ErrUomUnmarshall).WithError(err).With("id", id)
	}

	return &uom, nil
}

func (r *UomRepository) GetAll(ctx context.Context) ([]*domain.Uom, *zerrors.Error[domain.UomRepositoryError]) {
	docs, err := r.client.Collection(uomCollection).Documents(ctx).GetAll()
	if err != nil {
		return nil, zerrors.New(domain.ErrUomRead).WithError(err)
	}

	uoms := make([]*domain.Uom, 0, len(docs))
	for _, doc := range docs {
		var uom domain.Uom
		if err := doc.DataTo(&uom); err != nil {
			return nil, zerrors.New(domain.ErrUomUnmarshall).WithError(err).With("id", doc.Ref.ID)
		}
		uoms = append(uoms, &uom)
	}

	return uoms, nil
}

func (r *UomRepository) Delete(ctx context.Context, id string) *zerrors.Error[domain.UomRepositoryError] {
	_, err := r.client.Collection(uomCollection).Doc(id).Delete(ctx)
	if err != nil {
		return zerrors.New(domain.ErrUomDelete).WithError(err).With("id", id)
	}
	return nil
}
