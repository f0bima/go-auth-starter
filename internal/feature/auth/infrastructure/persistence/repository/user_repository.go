package repository

import (
	"context"
	"errors"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/entity"
	auth "github.com/f0bima/go-auth-starter/internal/feature/auth/domain/repository"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/infrastructure/persistence/mapper"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/infrastructure/persistence/schema"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Compile-time interface checks
var (
	_ auth.UserRepository         = (*userRepo)(nil)
	_ auth.RefreshTokenRepository = (*refreshTokenRepo)(nil)
)

type userRepo struct {
	db *gorm.DB
}

// NewUserRepository creates a new User repository.
func NewUserRepository(db *gorm.DB) auth.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(ctx context.Context, user *entity.User) error {
	s := mapper.UserEntityToSchema(user)
	err := r.db.WithContext(ctx).Create(s).Error
	if err != nil {
		return err
	}
	// Copy back generated ID and timestamps
	user.ID = s.ID
	user.CreatedAt = s.CreatedAt
	user.UpdatedAt = s.UpdatedAt
	return nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var s schema.UserSchema
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&s).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mapper.UserSchemaToEntity(&s), nil
}

func (r *userRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var s schema.UserSchema
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&s).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mapper.UserSchemaToEntity(&s), nil
}

type refreshTokenRepo struct {
	db *gorm.DB
}

// NewRefreshTokenRepository creates a new RefreshToken repository.
func NewRefreshTokenRepository(db *gorm.DB) auth.RefreshTokenRepository {
	return &refreshTokenRepo{db: db}
}

func (r *refreshTokenRepo) StoreRefreshToken(ctx context.Context, rt *entity.RefreshToken) error {
	s := mapper.RefreshTokenEntityToSchema(rt)
	err := r.db.WithContext(ctx).Create(s).Error
	if err != nil {
		return err
	}
	rt.ID = s.ID
	rt.CreatedAt = s.CreatedAt
	return nil
}

func (r *refreshTokenRepo) GetRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	var s schema.RefreshTokenSchema
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&s).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mapper.RefreshTokenSchemaToEntity(&s), nil
}

func (r *refreshTokenRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&schema.RefreshTokenSchema{}).Error
}
