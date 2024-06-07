package middleware

import (
	"net/http"
)

// CORSMiddleware is a middleware function that adds CORS headers to HTTP responses.
func ApplyCORSMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Call the next handler
		handler.ServeHTTP(w, r)
	})
}
