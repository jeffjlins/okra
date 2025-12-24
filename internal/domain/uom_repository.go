package domain

import "context"

// UomRepository defines the interface for Uom persistence operations
// This is a port in hexagonal architecture - implemented by adapters
type UomRepository interface {
	Save(ctx context.Context, uom *Uom) error
	GetByID(ctx context.Context, id string) (*Uom, error)
}
