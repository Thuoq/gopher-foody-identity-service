package auth

import (
	"gopher-identity-service/internal/presentation/http/middleware"
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
	authGroup.Use(middleware.GatewayAuth()) // Nhận diện user từ Gateway
	{
		authGroup.POST("/sign-up", r.signUpHandler.Handle)
		authGroup.POST("/sign-in", r.signInHandler.Handle)
		authGroup.POST("/refresh", r.refreshTokenHandler.Handle)

		// Protected routes (Xác thực đã được thực hiện tại Gateway)
		authGroup.POST("/logout", r.logoutHandler.Handle)
	}
}
