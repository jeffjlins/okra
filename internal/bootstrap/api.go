package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	httpadapter "github.com/jeffjlins/okra/internal/adapters/inbound/http"
	"github.com/jeffjlins/okra/internal/adapters/outbound/firestore"
	"github.com/jeffjlins/okra/internal/config"
)

type App struct {
	Server     *http.Server
	Firestore  *firestore.Client
	Repository *firestore.Repository
}

func NewApp(cfg *config.Config) (*App, error) {
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

	// Create repository
	repo := firestore.NewRepository(fsClient)

	// Create router with repository
	mux := httpadapter.NewRouter(repo)

	server := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		Server:     server,
		Firestore:  fsClient,
		Repository: repo,
	}, nil
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.Firestore.Close(); err != nil {
		return fmt.Errorf("failed to close firestore client: %w", err)
	}
	return a.Server.Shutdown(ctx)
}
