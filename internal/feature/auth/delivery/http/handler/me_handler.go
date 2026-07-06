package handler

import (
	"context"

	_ "github.com/f0bima/go-auth-starter/internal/feature/auth/delivery/http/dto"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/delivery/http/mapper"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/entity"
	"github.com/f0bima/go-core/response"
	"github.com/gin-gonic/gin"
)

// MeUseCase defines the use case interface for getting current user.
type MeUseCase interface {
	Me(ctx context.Context, userID string) (*entity.User, error)
}

type MeHandler struct {
	useCase MeUseCase
}

// NewMeHandler creates a new instance of MeHandler
func NewMeHandler(useCase MeUseCase) *MeHandler {
	return &MeHandler{useCase: useCase}
}

// Me @Summary Get current user
// @Description Get current user information based on JWT token
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=dto.UserResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /auth/me [get]
func (h *MeHandler) Handle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "UNAUTHORIZED", "Missing user_id in context")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		response.Unauthorized(c, "UNAUTHORIZED", "Invalid user_id in context")
		return
	}

	user, err := h.useCase.Me(c.Request.Context(), userIDStr)
	if err != nil {
		response.NotFound(c, "USER_NOT_FOUND", "User not found")
		return
	}

	// Assuming there is a ToUserResponse in mapper or we just return raw/DTO
	response.OK(c, mapper.ToUserResponse(user))
}
