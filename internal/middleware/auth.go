package middleware

import (
	"byte-board/internal/auth"
	"byte-board/internal/model"
	"context"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

// Stores user information in the request context
type contextKey string

const (
	UsernameContextKey contextKey = "username"
	RoleContextKey     contextKey = "role"
)

// Holds the JWT token provider for authentication
type AuthMiddleware struct {
	TokenProvider *auth.TokenProvider
}

// Creates a new authentication middleware
func NewAuthMiddleware(tokenProvider *auth.TokenProvider) *AuthMiddleware {
	return &AuthMiddleware{
		TokenProvider: tokenProvider,
	}
}

// Middleware that validates JWT tokens and adds user info to context
func (am *AuthMiddleware) JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract Authorization header
		authHeader := r.Header.Get("Authorization")

		// Check if authorization header exists
		if authHeader == "" {
			log.Warn().Msg("Missing authorization header")
			http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
			return
		}

		// Extract from "Bearer <token>" format"
		tokenString, err := extractBearerToken(authHeader)
		if err != nil {
			log.Warn().Err(err).Msg("Invalid Authorization header format")
			http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
			return
		}

		// Validate token
		if err := am.TokenProvider.ValidateToken(tokenString); err != nil {
			log.Warn().Err(err).Msg("Token validation failed")
			http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Parse the token to get claims
		claims, err := am.TokenProvider.ParseToken(tokenString)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to parse token claims")
			http.Error(w, "Unauthorized: Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Add username and role to request context
		ctx := context.WithValue(r.Context(), UsernameContextKey, claims.Username)
		ctx = context.WithValue(ctx, RoleContextKey, claims.Role)

		log.Debug().
			Str("username", claims.Username).
			Str("role", claims.Role).
			Str("path", r.URL.Path).
			Msg("User authenticated")

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Middleware that validates JWT if present, but allows requests without tokens
func (am *AuthMiddleware) OptionalJWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// If no auth header, just continue without adding user to context
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Try to extract and validate token
		tokenString, err := extractBearerToken(authHeader)
		if err != nil {
			// Invalid format, continue without auth
			next.ServeHTTP(w, r)
			return
		}

		// Validate token
		if err := am.TokenProvider.ValidateToken(tokenString); err != nil {
			// Invalid token, continue without auth
			next.ServeHTTP(w, r)
			return
		}

		// Parse token to get claims
		claims, err := am.TokenProvider.ParseToken(tokenString)
		if err != nil {
			// Can't parse claims, continue without auth
			next.ServeHTTP(w, r)
			return
		}

		// Add username and role to context
		ctx := context.WithValue(r.Context(), UsernameContextKey, claims.Username)
		ctx = context.WithValue(ctx, RoleContextKey, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Extracts the JWTtoken from "Bearer <token>" format
func extractBearerToken(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 {
		return "", model.ErrInvalidToken
	}
	if parts[0] != "Bearer" {
		return "", model.ErrInvalidToken
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", model.ErrInvalidToken
	}

	return token, nil
}

// Extracts username from the request context
func GetUsername(r *http.Request) string {
	username, ok := r.Context().Value(UsernameContextKey).(string) //<-- Type assertion
	if !ok {
		return ""
	}

	return username
}

// Extracts role from the request context
func GetRole(r *http.Request) string {
	role, ok := r.Context().Value(RoleContextKey).(string) //<-- Type assertion
	if !ok {
		return ""
	}

	return role
}

// Checks if authenticated user has a specific role
func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := GetRole(r)

			if role == "" {
				log.Warn().Msg("Role not found in context - ensure JWTAuth middleware is applied first")
				http.Error(w, "Forbidden: No role information", http.StatusForbidden)
				return
			}
			if role != requiredRole {
				log.Warn().
					Str("required_role", requiredRole).
					Str("user_role", role).
					Msg("User does not have required role")
				http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
