package handler

import (
	"context"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/delivery/http/dto"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/delivery/http/mapper"
	"github.com/f0bima/go-core/response"
	"github.com/gin-gonic/gin"
)

// RefreshUseCase defines the use case interface for token refresh.
type RefreshUseCase interface {
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
}

type RefreshHandler struct {
	useCase RefreshUseCase
}

// NewRefreshHandler creates a new instance of RefreshHandler
func NewRefreshHandler(useCase RefreshUseCase) *RefreshHandler {
	return &RefreshHandler{useCase: useCase}
}

// Refresh @Summary Token refresh
// @Description Get a new access token using a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshPayload true "Refresh request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/refresh [post]
func (h *RefreshHandler) Handle(c *gin.Context) {
	var req dto.RefreshPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	accessToken, newRefreshToken, err := h.useCase.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "UNAUTHORIZED", err.Error())
		return
	}

	response.OK(c, mapper.ToTokenResponse(accessToken, newRefreshToken))
}
