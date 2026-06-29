package entity

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the domain.
// This is a pure domain entity with no persistence concerns.
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// RefreshToken represents a refresh token in the domain.
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}
