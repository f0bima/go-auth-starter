package mapper

import (
	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/entity"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/infrastructure/persistence/schema"
)

// UserSchemaToEntity converts a GORM UserSchema to a domain User entity.
func UserSchemaToEntity(s *schema.UserSchema) *entity.User {
	if s == nil {
		return nil
	}
	return &entity.User{
		ID:           s.ID,
		Email:        s.Email,
		PasswordHash: s.PasswordHash,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}

// UserEntityToSchema converts a domain User entity to a GORM UserSchema.
func UserEntityToSchema(e *entity.User) *schema.UserSchema {
	if e == nil {
		return nil
	}
	return &schema.UserSchema{
		ID:           e.ID,
		Email:        e.Email,
		PasswordHash: e.PasswordHash,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

// RefreshTokenSchemaToEntity converts a GORM RefreshTokenSchema to a domain RefreshToken entity.
func RefreshTokenSchemaToEntity(s *schema.RefreshTokenSchema) *entity.RefreshToken {
	if s == nil {
		return nil
	}
	return &entity.RefreshToken{
		ID:        s.ID,
		UserID:    s.UserID,
		Token:     s.Token,
		ExpiresAt: s.ExpiresAt,
		CreatedAt: s.CreatedAt,
	}
}

// RefreshTokenEntityToSchema converts a domain RefreshToken entity to a GORM RefreshTokenSchema.
func RefreshTokenEntityToSchema(e *entity.RefreshToken) *schema.RefreshTokenSchema {
	if e == nil {
		return nil
	}
	return &schema.RefreshTokenSchema{
		ID:        e.ID,
		UserID:    e.UserID,
		Token:     e.Token,
		ExpiresAt: e.ExpiresAt,
		CreatedAt: e.CreatedAt,
	}
}
