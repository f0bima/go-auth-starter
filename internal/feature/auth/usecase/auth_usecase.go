package usecase

import (
	auth "github.com/f0bima/go-auth-starter/internal/feature/auth/domain/repository"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/service"
	authconfig "github.com/f0bima/go-auth-starter/internal/feature/auth/infrastructure/config"
	"go.opentelemetry.io/otel/trace"
)

type authUseCase struct {
	userRepo         auth.UserRepository
	refreshTokenRepo auth.RefreshTokenRepository
	cfg              *authconfig.AuthConfig
	tokenGen         service.TokenGenerator
	tracer           trace.Tracer
}

func NewAuthUseCase(
	userRepo auth.UserRepository,
	refreshTokenRepo auth.RefreshTokenRepository,
	cfg *authconfig.AuthConfig,
	tokenGen service.TokenGenerator,
	tracer trace.Tracer,
) *authUseCase {
	return &authUseCase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		cfg:              cfg,
		tokenGen:         tokenGen,
		tracer:           tracer,
	}
}
