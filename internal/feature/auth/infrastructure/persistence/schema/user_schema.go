package schema

import (
	"time"

	"github.com/google/uuid"
)

// UserSchema is the GORM model for users table.
type UserSchema struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM.
func (UserSchema) TableName() string {
	return "users"
}

// RefreshTokenSchema is the GORM model for refresh_tokens table.
type RefreshTokenSchema struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// TableName specifies the table name for GORM.
func (RefreshTokenSchema) TableName() string {
	return "refresh_tokens"
}
