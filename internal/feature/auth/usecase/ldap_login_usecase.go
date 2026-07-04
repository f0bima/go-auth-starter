package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/entity"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/repository"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/service"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LDAPLoginUseCase struct {
	ldapService      service.LDAPService
	userRepo         auth.UserRepository
	refreshTokenRepo auth.RefreshTokenRepository
	tokenGenerator   service.TokenGenerator
}

func NewLDAPLoginUseCase(
	ldapService service.LDAPService,
	userRepo auth.UserRepository,
	refreshTokenRepo auth.RefreshTokenRepository,
	tokenGenerator service.TokenGenerator,
) *LDAPLoginUseCase {
	return &LDAPLoginUseCase{
		ldapService:      ldapService,
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		tokenGenerator:   tokenGenerator,
	}
}

// LDAPLogin authenticates a user via LDAP and returns JWT tokens.
func (u *LDAPLoginUseCase) LDAPLogin(ctx context.Context, email, password string) (string, string, error) {
	// 1. Authenticate with LDAP
	valid, err := u.ldapService.Authenticate(ctx, email, password)
	if err != nil {
		return "", "", fmt.Errorf("ldap authentication failed: %w", err)
	}
	if !valid {
		return "", "", errors.New("invalid credentials")
	}

	// 2. Check if user exists in the local database
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("failed to check existing user: %w", err)
	}

	if user == nil {
		// Auto-register user if not found
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(uuid.New().String()), bcrypt.DefaultCost)
		if hashErr != nil {
			return "", "", fmt.Errorf("failed to hash password for auto-registration: %w", hashErr)
		}

		user = &entity.User{
			ID:           uuid.New(),
			Email:        email,
			PasswordHash: string(hashedPassword),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if createErr := u.userRepo.CreateUser(ctx, user); createErr != nil {
			return "", "", fmt.Errorf("failed to auto-register user: %w", createErr)
		}
	}

	// 3. Generate tokens
	accessToken, err := u.tokenGenerator.GenerateToken(user.ID.String(), user.Email, "access", 15*time.Minute)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := u.tokenGenerator.GenerateToken(user.ID.String(), user.Email, "refresh", 168*time.Hour)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 4. Store refresh token
	rtEntity := &entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(168 * time.Hour), // 7 days
		CreatedAt: time.Now(),
	}

	if err := u.refreshTokenRepo.StoreRefreshToken(ctx, rtEntity); err != nil {
		return "", "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}
