package usecase

import (
	"context"
	"time"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/entity"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (s *authUseCase) generateTokens(ctx context.Context, user *entity.User) (string, string, error) {
	ctx, span := s.tracer.Start(ctx, "AuthUseCase.generateTokens")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", user.ID.String()), attribute.String("user.email", user.Email))

	accessDuration, err := time.ParseDuration(s.cfg.JWTExpire)
	if err != nil {
		accessDuration = 15 * time.Minute
	}

	refreshDuration, err := time.ParseDuration(s.cfg.RefreshExpire)
	if err != nil {
		refreshDuration = 168 * time.Hour
	}

	accessToken, err := s.tokenGen.GenerateToken(user.ID.String(), user.Email, "access", accessDuration)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", "", err
	}

	refreshToken, err := s.tokenGen.GenerateToken(user.ID.String(), user.Email, "refresh", refreshDuration)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", "", err
	}

	rt := &entity.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(refreshDuration),
	}
	err = s.refreshTokenRepo.StoreRefreshToken(ctx, rt)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
