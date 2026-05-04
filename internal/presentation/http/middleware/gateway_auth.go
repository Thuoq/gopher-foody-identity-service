package middleware

import (
	"github.com/gin-gonic/gin"
)

// GatewayAuth extracts user identity from headers injected by the API Gateway.
func GatewayAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-Id")
		sessionID := c.GetHeader("X-Session-Id")

		if userID != "" {
			c.Set("public_user_id", userID)
		}
		
		if sessionID != "" {
			c.Set("session_id", sessionID)
		}

		c.Next()
	}
}
