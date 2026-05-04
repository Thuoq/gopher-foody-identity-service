package auth

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	signUpHandler       *SignUpHandler
	signInHandler       *SignInHandler
	refreshTokenHandler *RefreshTokenHandler
}

func NewRouter(
	signUpHandler *SignUpHandler,
	signInHandler *SignInHandler,
	refreshTokenHandler *RefreshTokenHandler,
) *Router {
	return &Router{
		signUpHandler:       signUpHandler,
		signInHandler:       signInHandler,
		refreshTokenHandler: refreshTokenHandler,
	}
}

func (r *Router) Register(api *gin.RouterGroup) {
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/sign-up", r.signUpHandler.Handle)
		authGroup.POST("/sign-in", r.signInHandler.Handle)
		authGroup.POST("/refresh", r.refreshTokenHandler.Handle)
	}
}
