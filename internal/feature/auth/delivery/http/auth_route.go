package http

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all auth-specific routes on the given router.
func (h *AuthController) RegisterRoutes(r *gin.Engine) {
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", h.Register)
		authRoutes.POST("/login", h.Login)
		authRoutes.POST("/refresh", h.Refresh)
	}

	r.GET("/.well-known/jwks.json", h.JWKS)
}
