package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jeffjlins/okra/internal/adapters/outbound/firestore"
)

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := map[string]string{
		"message": "hello, world",
	}

	json.NewEncoder(w).Encode(resp)
}

// DemoRequest represents the JSON payload for the demo endpoint
type DemoRequest struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
}

func createDemoHandler(repo *firestore.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// Decode JSON payload
		var req DemoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		// Validate that both fields are provided
		if req.Field1 == "" || req.Field2 == "" {
			http.Error(w, "Both field1 and field2 are required", http.StatusBadRequest)
			return
		}

		// Save to Firestore
		ctx := r.Context()
		docID, err := repo.SaveWithAutoID(ctx, "demo", req)
		if err != nil {
			log.Printf("Error saving to Firestore: %v", err)
			http.Error(w, "Failed to save data", http.StatusInternalServerError)
			return
		}

		// Return success response with document ID
		resp := map[string]interface{}{
			"id":     docID,
			"field1": req.Field1,
			"field2": req.Field2,
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}
