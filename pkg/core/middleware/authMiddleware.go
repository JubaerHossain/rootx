package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/JubaerHossain/rootx/pkg/auth"
	"github.com/JubaerHossain/rootx/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

type key int

const (
	claimsKey key = iota
)

// AuthMiddleware validates JWT tokens and protects routes
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteJSONError(w, http.StatusUnauthorized, "Unauthorized: missing token")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.WriteJSONError(w, http.StatusUnauthorized, "Invalid Authorization Header Format")
			return
		}

		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			utils.WriteJSONError(w, http.StatusUnauthorized, fmt.Sprintf("Invalid Token: %v", err))
			return
		}

		// Add claims to context for use in other handlers
		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetClaimsFromContext retrieves claims from request context
func GetClaimsFromContext(ctx context.Context) (jwt.MapClaims, bool) {
	claims, ok := ctx.Value(claimsKey).(jwt.MapClaims)
	return claims, ok
}
