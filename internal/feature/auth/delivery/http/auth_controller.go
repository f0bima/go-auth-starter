package http

import (
	"net/http"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/delivery/http/dto"
	auth "github.com/f0bima/go-auth-starter/internal/feature/auth/domain/repository"
	"github.com/f0bima/go-core/response"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	useCase auth.AuthUseCase
}

// NewAuthController creates a new instance of AuthController
func NewAuthController(useCase auth.AuthUseCase) *AuthController {
	return &AuthController{useCase: useCase}
}

// Register @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.AuthPayload true "Register request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *AuthController) Register(c *gin.Context) {
	var req dto.AuthPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "BAD_REQUEST", err.Error())
		return
	}

	user, err := h.useCase.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.BadRequest(c, "BAD_REQUEST", err.Error())
		return
	}

	response.Created(c, ToUserResponse(user))
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
func (h *AuthController) Login(c *gin.Context) {
	var req dto.AuthPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "BAD_REQUEST", err.Error())
		return
	}

	accessToken, refreshToken, err := h.useCase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Unauthorized(c, "UNAUTHORIZED", err.Error())
		return
	}

	response.OK(c, ToTokenResponse(accessToken, refreshToken))
}

// Refresh @Summary Refresh token
// @Description Get a new access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshPayload true "Refresh token request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/refresh [post]
func (h *AuthController) Refresh(c *gin.Context) {
	var req dto.RefreshPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "BAD_REQUEST", err.Error())
		return
	}

	accessToken, refreshToken, err := h.useCase.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "UNAUTHORIZED", err.Error())
		return
	}

	response.OK(c, ToTokenResponse(accessToken, refreshToken))
}

// JWKS @Summary Get JWKS
// @Description Get the JSON Web Key Set for token validation by external services
// @Tags auth
// @Produce json
// @Success 200 {object} dto.JWKSResponse "JWKS"
// @Router /.well-known/jwks.json [get]
func (h *AuthController) JWKS(c *gin.Context) {
	c.JSON(http.StatusOK, ToJWKSResponse(h.useCase.GetJWKS()))
}
