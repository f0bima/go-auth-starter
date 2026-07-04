package usecase

import (
	"context"
	"errors"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/entity"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (s *authUseCase) Me(ctx context.Context, userID string) (*entity.User, error) {
	ctx, span := s.tracer.Start(ctx, "AuthUseCase.Me")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", userID))

	uid, err := uuid.Parse(userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetUserByID(ctx, uid)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	if user == nil {
		err = errors.New("user not found")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return user, nil
}
