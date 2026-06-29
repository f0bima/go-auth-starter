package dto

import (
	"github.com/google/uuid"
)

// UserResponse is the HTTP response for user registration.
type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

// TokenResponse is the HTTP response for login and token refresh.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// JWKResponse represents a single JSON Web Key in the response.
type JWKResponse struct {
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// JWKSResponse is the HTTP response for the JWKS endpoint.
type JWKSResponse struct {
	Keys []JWKResponse `json:"keys"`
}
