package usecase

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (s *authUseCase) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	ctx, span := s.tracer.Start(ctx, "AuthUseCase.Refresh")
	defer span.End()

	claims, err := s.tokenGen.ValidateToken(refreshToken)
	if err != nil {
		err = errors.New("invalid refresh token")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", "", err
	}

	span.SetAttributes(attribute.String("user.id", claims.UserID))

	if claims.Type != "refresh" {
		err = errors.New("invalid token type")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", "", err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		err = errors.New("invalid user ID in token")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", "", err
	}

	rt, err := s.refreshTokenRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", "", err
	}
	if rt == nil {
		err = errors.New("refresh token not found or already revoked")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", "", err
	}

	if rt.ExpiresAt.Before(time.Now()) {
		_ = s.refreshTokenRepo.DeleteRefreshToken(ctx, refreshToken)
		err = errors.New("refresh token expired")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", "", err
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", "", err
	}
	if user == nil {
		err = errors.New("user not found")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", "", err
	}

	_ = s.refreshTokenRepo.DeleteRefreshToken(ctx, refreshToken)

	slog.InfoContext(ctx, "Token successfully refreshed", "event", "token_refreshed", "userId", userID.String())

	return s.generateTokens(ctx, user)
}
