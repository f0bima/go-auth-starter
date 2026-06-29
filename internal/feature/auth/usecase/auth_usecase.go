package usecase

import (
	auth "github.com/f0bima/go-auth-starter/internal/feature/auth/domain/repository"
	"go.opentelemetry.io/otel/trace"
)

type authUseCase struct {
	userRepo         auth.UserRepository
	refreshTokenRepo auth.RefreshTokenRepository
	cfg              *AuthConfig
	tokenGen         auth.TokenGenerator
	tracer           trace.Tracer
}

func NewAuthUseCase(
	userRepo auth.UserRepository,
	refreshTokenRepo auth.RefreshTokenRepository,
	cfg *AuthConfig,
	tokenGen auth.TokenGenerator,
	tracer trace.Tracer,
) auth.AuthUseCase {
	return &authUseCase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		cfg:              cfg,
		tokenGen:         tokenGen,
		tracer:           tracer,
	}
}
