package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

type Repository struct {
	client *Client
}

func NewRepository(client *Client) *Repository {
	return &Repository{
		client: client,
	}
}

// Save saves a document to the specified collection
func (r *Repository) Save(ctx context.Context, collection string, id string, data interface{}) error {
	_, err := r.client.Collection(collection).Doc(id).Set(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to save document to %s/%s: %w", collection, id, err)
	}
	return nil
}

// SaveWithAutoID saves a document to the specified collection with an auto-generated ID
// Returns the generated document ID
func (r *Repository) SaveWithAutoID(ctx context.Context, collection string, data interface{}) (string, error) {
	docRef := r.client.Collection(collection).NewDoc()
	_, err := docRef.Set(ctx, data)
	if err != nil {
		return "", fmt.Errorf("failed to save document to %s: %w", collection, err)
	}
	return docRef.ID, nil
}

// Get retrieves a document from the specified collection
func (r *Repository) Get(ctx context.Context, collection string, id string) (*firestore.DocumentSnapshot, error) {
	doc, err := r.client.Collection(collection).Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get document from %s/%s: %w", collection, id, err)
	}
	return doc, nil
}

// Delete removes a document from the specified collection
func (r *Repository) Delete(ctx context.Context, collection string, id string) error {
	_, err := r.client.Collection(collection).Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete document from %s/%s: %w", collection, id, err)
	}
	return nil
}

// GetAll retrieves all documents from the specified collection
func (r *Repository) GetAll(ctx context.Context, collection string) ([]*firestore.DocumentSnapshot, error) {
	docs, err := r.client.Collection(collection).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all documents from %s: %w", collection, err)
	}
	return docs, nil
}
