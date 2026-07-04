package auth

import (
	"context"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/entity"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user persistence operations.
type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
}
