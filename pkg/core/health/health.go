// Package health_check provides functionality for health checks.
package health

import (
	"encoding/json"
	"net/http"

	"github.com/JubaerHossain/rootx/pkg/core/app"
	"go.uber.org/zap"
)

// HealthCheckResponse represents the structure of the health check response
type HealthCheckResponse struct {
	Status string `json:"status"`
}

// HealthCheckHandler returns an HTTP handler function for health checks
func HealthCheckHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create the health check response
		response := HealthCheckResponse{Status: "OK"}

		// Set the response header to indicate JSON content
		w.Header().Set("Content-Type", "application/json")

		// Write the response body
		if err := json.NewEncoder(w).Encode(response); err != nil {
			// Log the error if writing the response fails
			app.Logger.Error("Error writing response", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Response was successfully written
		w.WriteHeader(http.StatusOK)
	}
}

func DBHealthHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check the health of the database
		if err := app.CheckDatabaseHealth(); err != nil {
			app.Logger.Error("Database health check failed", zap.Error(err))
			http.Error(w, "Database is unhealthy", http.StatusServiceUnavailable)
		}

		// If the database is healthy
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Database is healthy"))
	}
}
