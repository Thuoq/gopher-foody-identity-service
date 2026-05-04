package middleware

import (
	"gopher-identity-service/pkg/jwt"
	"gopher-identity-service/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthGuard(jwtManager jwt.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "authorization header must be Bearer token")
			c.Abort()
			return
		}

		accessToken := parts[1]
		claims, err := jwtManager.ValidateAccessToken(accessToken)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "invalid or expired access token")
			c.Abort()
			return
		}

		// Inject into context
		c.Set("public_user_id", claims.PublicUserId)
		c.Set("session_id", claims.SessionID)

		c.Next()
	}
}
