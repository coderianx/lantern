package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/coderianx/lantern/internal/server/models"
	"github.com/golang-jwt/jwt/v5"
)

// claimsKey is the private context key used to store *models.Claims.
// Unexported to avoid collisions with other packages.
type claimsKey struct{}

// Auth is a chi-compatible middleware.
//
// Usage:
//
//	r.With(middlewares.Auth).Get("/buckets", handlers.CreateNewBucket)
//	r.Group(func(r chi.Router) {
//		r.Use(middlewares.Auth)
//		...
//	})
//
// It expects: Authorization: Bearer <accessToken>
// On success it stores *models.Claims in the request context and calls next.
// On failure it aborts with 401 (or 500 when JWT_SECRET is missing).
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		tokenStr := tokenFromHeader(r.Header.Get("Authorization"))
		if tokenStr == "" {
			writeError(w, http.StatusUnauthorized, "missing or invalid Authorization header")
			return
		}

		secret := os.Getenv("JWT_SECRET")
		if len(secret) == 0 {
			writeError(w, http.StatusInternalServerError, "JWT_SECRET is not configured")
			return
		}

		claims := &models.Claims{}

		token, err := jwt.ParseWithClaims(
			tokenStr,
			claims,
			func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			},
		)
		if err != nil || !token.Valid {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey{}, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ClaimsFromContext returns the *models.Claims stored by Auth.
// Intended use inside handlers:
//
//	claims, ok := middlewares.ClaimsFromContext(r.Context())
//	if !ok {
//		// Auth middleware missing or bypassed
//	}
func ClaimsFromContext(ctx context.Context) (*models.Claims, bool) {
	claims, ok := ctx.Value(claimsKey{}).(*models.Claims)
	if !ok || claims == nil {
		return nil, false
	}
	return claims, true
}

// MustClaims is like ClaimsFromContext but returns nil when absent.
// Useful when the handler is guaranteed to run behind Auth.
func MustClaims(ctx context.Context) *models.Claims {
	claims, _ := ClaimsFromContext(ctx)
	return claims
}

// tokenFromHeader extracts the token from "Bearer <token>".
// Returns "" when the header is missing or malformed.
func tokenFromHeader(header string) string {
	if header == "" {
		return ""
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return ""
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return strings.TrimSpace(parts[1])
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": msg,
	})
}
