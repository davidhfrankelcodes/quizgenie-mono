// backend/internal/auth/middleware.go
package auth

import (
	"context"
	"net/http"
	"strings"
)

// key type so context is not conflicted
type contextKey string

const userContextKey = contextKey("userClaims")

// AuthMiddleware ensures the request has a valid JWT. It sets the Claims
// in request.Context, so downstream handlers can read it.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Store claims in context for later handlers
		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FromContext retrieves the JWT claims that were stored by AuthMiddleware.
func FromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(userContextKey).(*Claims)
	return claims, ok
}
