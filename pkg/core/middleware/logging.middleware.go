// File: wild-workouts-go-ddd-example/pkg/core/middleware/logging.go

package middleware

import (
	"net/http"
	"time"

	"github.com/JubaerHossain/rootx/pkg/core/logger"
	"go.uber.org/zap"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Call the next handler
		next.ServeHTTP(w, r)

		// Log request details after handling
		logger.Info("Request handled",
			zap.String("method", r.Method),
			zap.String("url", r.URL.Path),
			zap.Duration("duration", time.Since(start)),
		)
	})
}
