package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

type Client struct {
	*firestore.Client
}

// NewClient creates a Firestore client with default credentials
// Uses GOOGLE_APPLICATION_CREDENTIALS env var or default credentials if available
func NewClient(ctx context.Context, projectID string, databaseID string) (*Client, error) {
	client, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %w", err)
	}

	return &Client{Client: client}, nil
}

// NewClientWithCredentials creates a Firestore client using a service account JSON file
func NewClientWithCredentials(ctx context.Context, projectID string, databaseID string, credentialsFile string) (*Client, error) {
	opts := []option.ClientOption{
		option.WithCredentialsFile(credentialsFile),
	}

	client, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %w", err)
	}

	return &Client{Client: client}, nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}
