package route

import (
	"crypto/rsa"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/delivery/http/handler"
	"github.com/f0bima/go-core/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the routing for auth endpoints.
func RegisterRoutes(
	r *gin.Engine,
	pubKey *rsa.PublicKey,
	registerHandler *handler.RegisterHandler,
	loginHandler *handler.LoginHandler,
	ldapLoginHandler *handler.LDAPLoginHandler,
	refreshHandler *handler.RefreshHandler,
	jwksHandler *handler.JwksHandler,
	meHandler *handler.MeHandler,
	panicHandler *handler.PanicHandler,
) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", registerHandler.Handle)
		authGroup.POST("/login", loginHandler.Handle)
		authGroup.POST("/login/ldap", ldapLoginHandler.Handle)
		authGroup.POST("/refresh", refreshHandler.Handle)

		// Protected endpoints
		protectedGroup := authGroup.Group("")
		protectedGroup.Use(middleware.Auth(pubKey))
		{
			protectedGroup.GET("/me", meHandler.Handle)
		}
	}

	// JWKS endpoint
	r.GET("/.well-known/jwks.json", jwksHandler.Handle)

	// Test panic
	authGroup.GET("/panic", panicHandler.Handle)
}
