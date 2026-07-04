package handler

import (
	"net/http"

	_ "github.com/f0bima/go-auth-starter/internal/feature/auth/delivery/http/dto"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/delivery/http/mapper"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/service"
	"github.com/gin-gonic/gin"
)

// JwksUseCase defines the use case interface for fetching JWKS.
type JwksUseCase interface {
	GetJWKS() service.JWKS
}

type JwksHandler struct {
	useCase JwksUseCase
}

// NewJwksHandler creates a new instance of JwksHandler
func NewJwksHandler(useCase JwksUseCase) *JwksHandler {
	return &JwksHandler{useCase: useCase}
}

// JWKS @Summary Get JWKS
// @Description Get the JSON Web Key Set for token validation by external services
// @Tags auth
// @Produce json
// @Success 200 {object} dto.JWKSResponse "JWKS"
// @Router /.well-known/jwks.json [get]
func (h *JwksHandler) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, mapper.ToJWKSResponse(h.useCase.GetJWKS()))
}
