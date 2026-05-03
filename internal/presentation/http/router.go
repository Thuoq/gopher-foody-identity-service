package http

import (
	"gopher-identity-service/internal/presentation/http/handlers/auth"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"gopher-identity-service/internal/config"
	"gopher-identity-service/internal/presentation/http/handlers/user"
)

func NewRouter(cfg *config.Config, logger *zap.Logger, userRouter *user.Router, authRouter *auth.Router) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Configure gin validator to use json tags instead of struct field names
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	// Add basic middlewares
	r.Use(gin.Recovery())

	// Example health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	api := r.Group("/api/v1")
	userRouter.Register(api)
	authRouter.Register(api)

	return r
}
