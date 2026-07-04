package service

import (
	"time"
)

// TokenClaims represents the decoded claims from a token.
type TokenClaims struct {
	UserID string
	Email  string
	Type   string // "access" or "refresh"
}

// JWKS represents the JSON Web Key Set for public key distribution.
type JWKS struct {
	Keys []JWK
}

// JWK represents a single JSON Web Key.
type JWK struct {
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// TokenGenerator defines the interface for token generation and validation.
type TokenGenerator interface {
	GenerateToken(userID string, email string, tokenType string, expiration time.Duration) (string, error)
	ValidateToken(tokenString string) (*TokenClaims, error)
	GetJWKS() JWKS
}
