package middleware

import (
	"github.com/gin-gonic/gin"
)

// AuthRequiredMiddleware validates JWT access tokens from the Authorization header.
// This is an auth-specific middleware for protecting routes that require authentication.
func AuthRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement JWT validation
		// 1. Extract Bearer token from Authorization header
		// 2. Validate token using RSA public key
		// 3. Set user info in gin.Context
		c.Next()
	}
}
