package auth

import (
	"context"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/entity"
)

// RefreshTokenRepository defines the interface for refresh token persistence operations.
type RefreshTokenRepository interface {
	StoreRefreshToken(ctx context.Context, rt *entity.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}
