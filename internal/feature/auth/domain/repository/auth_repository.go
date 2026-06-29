package auth

import (
	"context"
	"time"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/entity"
	"github.com/google/uuid"
)

// TokenGenerator defines the interface for token generation and validation.
// Implemented by infrastructure/jwt, keeping domain independent of infrastructure.
type TokenGenerator interface {
	GenerateToken(userID string, email string, tokenType string, expiration time.Duration) (string, error)
	ValidateToken(tokenString string) (*TokenClaims, error)
	GetJWKS() JWKS
}

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

// UserRepository defines the interface for user persistence operations.
type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
}

// RefreshTokenRepository defines the interface for refresh token persistence operations.
type RefreshTokenRepository interface {
	StoreRefreshToken(ctx context.Context, rt *entity.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}

// OutboxRepository defines the interface for outbox event persistence.
type OutboxRepository interface {
	SaveEvent(ctx context.Context, event *entity.OutboxEvent) error
	GetPendingEvents(ctx context.Context, limit int) ([]*entity.OutboxEvent, error)
	GetRetryableEvents(ctx context.Context, limit int) ([]*entity.OutboxEvent, error)
	MarkEventProcessed(ctx context.Context, id uuid.UUID) error
	MarkEventFailed(ctx context.Context, id uuid.UUID, err error) error
}

// InboxRepository defines the interface for inbox event persistence (idempotent consumer).
type InboxRepository interface {
	SaveEventIfNotExists(ctx context.Context, event *entity.InboxEvent) (bool, error)
	GetPendingEvents(ctx context.Context, limit int) ([]*entity.InboxEvent, error)
	GetRetryableEvents(ctx context.Context, limit int) ([]*entity.InboxEvent, error)
	MarkEventProcessing(ctx context.Context, id uuid.UUID) error
	MarkEventProcessed(ctx context.Context, id uuid.UUID) error
	MarkEventFailed(ctx context.Context, id uuid.UUID, err error) error
}

// AuthUseCase defines the interface for authentication business logic.
type AuthUseCase interface {
	Register(ctx context.Context, email, password string) (*entity.User, error)
	Login(ctx context.Context, email, password string) (string, string, error) // returns (accessToken, refreshToken, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)  // returns (accessToken, refreshToken, error)
	GetJWKS() JWKS
}
