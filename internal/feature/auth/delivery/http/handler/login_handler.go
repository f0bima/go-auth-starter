package handler

import (
	"context"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/delivery/http/dto"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/delivery/http/mapper"
	"github.com/f0bima/go-core/response"
	"github.com/gin-gonic/gin"
)

// LoginUseCase defines the use case interface for logging in.
type LoginUseCase interface {
	Login(ctx context.Context, email, password string) (string, string, error)
}

type LoginHandler struct {
	useCase LoginUseCase
}

// NewLoginHandler creates a new instance of LoginHandler
func NewLoginHandler(useCase LoginUseCase) *LoginHandler {
	return &LoginHandler{useCase: useCase}
}

// Login @Summary User login
// @Description Login with email and password to get access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.AuthPayload true "Login request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *LoginHandler) Handle(c *gin.Context) {
	var req dto.AuthPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	accessToken, refreshToken, err := h.useCase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Unauthorized(c, "UNAUTHORIZED", err.Error())
		return
	}

	response.OK(c, mapper.ToTokenResponse(accessToken, refreshToken))
}
