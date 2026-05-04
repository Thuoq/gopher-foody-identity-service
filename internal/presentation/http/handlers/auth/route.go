package auth

import (
	"gopher-identity-service/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type Router struct {
	signUpHandler       *SignUpHandler
	signInHandler       *SignInHandler
	refreshTokenHandler *RefreshTokenHandler
	logoutHandler       *LogoutHandler
	jwtManager          jwt.TokenManager
}

func NewRouter(
	signUpHandler *SignUpHandler,
	signInHandler *SignInHandler,
	refreshTokenHandler *RefreshTokenHandler,
	logoutHandler *LogoutHandler,
	jwtManager jwt.TokenManager,
) *Router {
	return &Router{
		signUpHandler:       signUpHandler,
		signInHandler:       signInHandler,
		refreshTokenHandler: refreshTokenHandler,
		logoutHandler:       logoutHandler,
		jwtManager:          jwtManager,
	}
}

func (r *Router) Register(api *gin.RouterGroup) {
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/sign-up", r.signUpHandler.Handle)
		authGroup.POST("/sign-in", r.signInHandler.Handle)
		authGroup.POST("/refresh", r.refreshTokenHandler.Handle)

		// Protected routes
		authGroup.POST("/logout", r.logoutHandler.Handle)
	}
}
