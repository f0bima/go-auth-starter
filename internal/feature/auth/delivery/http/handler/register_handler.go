package handler

import (
	"context"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/delivery/http/dto"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/delivery/http/mapper"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/entity"
	"github.com/f0bima/go-core/response"
	"github.com/gin-gonic/gin"
)

// RegisterUseCase defines the use case interface for registering.
type RegisterUseCase interface {
	Register(ctx context.Context, email, password string) (*entity.User, error)
}

type RegisterHandler struct {
	useCase RegisterUseCase
}

// NewRegisterHandler creates a new instance of RegisterHandler
func NewRegisterHandler(useCase RegisterUseCase) *RegisterHandler {
	return &RegisterHandler{useCase: useCase}
}

// Register @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.AuthPayload true "Register request"
// @Success 201 {object} response.SuccessResponse{data=dto.UserResponse}
// @Failure 400 {object} response.ErrorResponse
// @Router /auth/register [post]
func (h *RegisterHandler) Handle(c *gin.Context) {
	var req dto.AuthPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	user, err := h.useCase.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.BadRequest(c, "BAD_REQUEST", err.Error())
		return
	}

	response.Created(c, mapper.ToUserResponse(user))
}
