package domain

import "context"

type UomRepository interface {
	Save(ctx context.Context, uom *Uom) error
	GetByID(ctx context.Context, id string) (*Uom, error)
	GetAll(ctx context.Context) ([]*Uom, error)
	Delete(ctx context.Context, id string) error
}
