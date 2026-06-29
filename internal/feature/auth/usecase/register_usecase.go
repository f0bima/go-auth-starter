package usecase

import (
	"context"
	"errors"
	"log/slog"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/entity"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/crypto/bcrypt"
)

func (s *authUseCase) Register(ctx context.Context, email, password string) (*entity.User, error) {
	ctx, span := s.tracer.Start(ctx, "AuthUseCase.Register")
	defer span.End()

	span.SetAttributes(attribute.String("user.email", email))

	existingUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	if existingUser != nil {
		err = errors.New("email already in use")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	user := &entity.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	slog.InfoContext(ctx, "User successfully registered", "event", "user_registered", "userId", user.ID.String())

	return user, nil
}
