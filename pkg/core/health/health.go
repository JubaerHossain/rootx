// Package health_check provides functionality for health checks.
package health

import (
	"log"
	"net/http"
)

// HealthCheckHandler returns an HTTP handler function for health checks.
func HealthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Write the response body
		_, err := w.Write([]byte("OK"))
		if err != nil {
			// If an error occurs while writing the response, log it
			log.Printf("Error writing response: %v", err)
			// You may choose to return an error response to the client, or handle it differently based on your requirements
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// Response was successfully written
		w.WriteHeader(http.StatusOK)
	}
}
