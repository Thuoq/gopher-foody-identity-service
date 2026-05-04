package auth

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	signUpHandler *SignUpHandler
	signInHandler *SignInHandler
}

func NewRouter(signUpHandler *SignUpHandler, signInHandler *SignInHandler) *Router {
	return &Router{
		signUpHandler: signUpHandler,
		signInHandler: signInHandler,
	}
}

func (r *Router) Register(api *gin.RouterGroup) {
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/sign-up", r.signUpHandler.Handle)
		authGroup.POST("/sign-in", r.signInHandler.Handle)
	}
}
