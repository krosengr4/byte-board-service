package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrInvalidSignature = errors.New("invalid token signature")
	ErrMissingClaims    = errors.New("missing required claims")
)

// JWT claims structure
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// JWT configuration
type JWTConfig struct {
	SecretKey       string
	ExpirationHours int
}

// JWT Token creation and validation
type TokenProvider struct {
	config JWTConfig
}

// Creates a new JWT token provider
func NewTokenProvider(config JWTConfig) *TokenProvider {
	return &TokenProvider{
		config: config,
	}
}

// Generates new JWT token for a given user
func (tp *TokenProvider) CreateToken(username string, role string) (string, error) {
	now := time.Now()
	expirationTime := now.Add(time.Duration(tp.config.ExpirationHours) * time.Hour)

	// Create claims with user info and standard class
	claims := &Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// Sign token with secret key (using HMAC-SHA512)
	tokenString, err := token.SignedString([]byte(tp.config.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// Validate the JWT token signature and expiration
func (tp *TokenProvider) ValidateToken(tokenString string) error {
	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{},
		error) {
		// Verify signing method is HMAC-SHA512
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tp.config.SecretKey), nil
	})

	if err != nil {
		// Check for specific JWT errors
		if errors.Is(err, jwt.ErrTokenExpired) {
			return ErrExpiredToken
		}
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return jwt.ErrSignatureInvalid
		}
		return fmt.Errorf("%w, %v", ErrInvalidToken, err)
	}

	// Verify that the token is valid
	if !token.Valid {
		return ErrInvalidToken
	}

	return nil
}

// Parse the JWT token and return the claims
func (tp *TokenProvider) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tp.config.SecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse the token: %w", err)
	}

	// Extract Claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Varify required claims exists
	if claims.Username == "" {
		return nil, ErrMissingClaims
	}

	return claims, nil
}

// Extract role from a JWT token
func (tp *TokenProvider) GetAuthoritiesFromToken(tokenString string) (string, error) {
	claims, err := tp.ParseToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("error parsing token: %w", err)
	}

	return claims.Role, nil
}

// Extract the JWT token from auth header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	const bearerPrefix = "Bearer "

	if len(authHeader) < len(bearerPrefix) {
		return "", errors.New("invalid authoirization header format")
	}
	if authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("authorization header must start with 'Bearer '")
	}

	token := authHeader[len(bearerPrefix):]
	if token == "" {
		return "", errors.New("token is empty")
	}

	return token, nil
}
