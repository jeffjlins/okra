package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	httpadapter "github.com/jeffjlins/okra/internal/adapters/inbound/http"
	"github.com/jeffjlins/okra/internal/adapters/outbound/firestore"
	"github.com/jeffjlins/okra/internal/usecase"
)

type App struct {
	Server     *http.Server
	Firestore  *firestore.Client
}

func NewApp(cfg *Config) (*App, error) {
	ctx := context.Background()

	// Initialize Firestore client
	var fsClient *firestore.Client
	var err error
	if cfg.Firestore.CredentialsFile != "" {
		fsClient, err = firestore.NewClientWithCredentials(ctx, cfg.Firestore.ProjectID, cfg.Firestore.DatabaseID, cfg.Firestore.CredentialsFile)
	} else {
		fsClient, err = firestore.NewClient(ctx, cfg.Firestore.ProjectID, cfg.Firestore.DatabaseID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to initialize firestore client: %w", err)
	}

	// Create repositories
	uomRepo := firestore.NewUomRepository(fsClient)

	// Create use cases/services
	uomService := usecase.NewUomService(uomRepo)

	// Create router with repositories and services
	mux := httpadapter.NewRouter(uomService)

	server := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		Server:     server,
		Firestore:  fsClient,
	}, nil
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.Firestore.Close(); err != nil {
		return fmt.Errorf("failed to close firestore client: %w", err)
	}
	return a.Server.Shutdown(ctx)
}
