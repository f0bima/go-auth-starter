package usecase

import (
	"context"
	"errors"
	"log/slog"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/crypto/bcrypt"
)

func (s *authUseCase) Login(ctx context.Context, email, password string) (string, string, error) {
	ctx, span := s.tracer.Start(ctx, "AuthUseCase.Login")
	defer span.End()

	span.SetAttributes(attribute.String("user.email", email))

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", "", err
	}
	if user == nil {
		err = errors.New("invalid credentials")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		err = errors.New("invalid credentials")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", "", err
	}

	slog.InfoContext(ctx, "User successfully logged in", "event", "user_logged_in", "email", email)

	return s.generateTokens(ctx, user)
}
