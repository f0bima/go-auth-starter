package handler

import (
	"context"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/delivery/http/dto"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/delivery/http/mapper"
	"github.com/f0bima/go-core/response"
	"github.com/gin-gonic/gin"
)

// LDAPLoginUseCase defines the use case interface for LDAP logging in.
type LDAPLoginUseCase interface {
	LDAPLogin(ctx context.Context, email, password string) (string, string, error)
}

type LDAPLoginHandler struct {
	useCase LDAPLoginUseCase
}

// NewLDAPLoginHandler creates a new instance of LDAPLoginHandler
func NewLDAPLoginHandler(useCase LDAPLoginUseCase) *LDAPLoginHandler {
	return &LDAPLoginHandler{useCase: useCase}
}

// Handle @Summary LDAP User login
// @Description Login with email/username and password via LDAP to get access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.AuthPayload true "LDAP Login request"
// @Success 200 {object} response.SuccessResponse{data=dto.TokenResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/login/ldap [post]
func (h *LDAPLoginHandler) Handle(c *gin.Context) {
	var req dto.LDAPLoginPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	accessToken, refreshToken, err := h.useCase.LDAPLogin(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Unauthorized(c, "UNAUTHORIZED", err.Error())
		return
	}

	response.OK(c, mapper.ToTokenResponse(accessToken, refreshToken))
}
