package auth

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	signUpHandler *SignUpHandler
}

func NewRouter(signUpHandler *SignUpHandler) *Router {
	return &Router{
		signUpHandler: signUpHandler,
	}
}

func (r *Router) Register(api *gin.RouterGroup) {
	userGroup := api.Group("/auth")
	{
		userGroup.GET("/signUp", r.signUpHandler.Handle)
	}
}
