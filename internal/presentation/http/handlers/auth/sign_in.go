package auth

import (
	"errors"
	"net/http"

	"gopher-identity-service/internal/application/usecases"
	"gopher-identity-service/internal/config"
	"gopher-identity-service/internal/core/ports"
	"gopher-identity-service/internal/presentation/http/handlers/auth/dto"
	"gopher-identity-service/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SignInHandler struct {
	signInUseCase ports.ISignInUseCase
	cfg           *config.Config
	logger        *zap.Logger
}

func NewSignInHandler(signInUseCase ports.ISignInUseCase, cfg *config.Config, logger *zap.Logger) *SignInHandler {
	return &SignInHandler{
		signInUseCase: signInUseCase,
		cfg:           cfg,
		logger:        logger,
	}
}

func (h *SignInHandler) Handle(c *gin.Context) {
	var body dto.SignInRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		fieldErrors := response.ParseValidationErrors(err)
		if fieldErrors != nil {
			response.ValidationError(c, http.StatusBadRequest, "invalid input data", fieldErrors)
			return
		}

		h.logger.Error("invalid json body", zap.Error(err))
		response.Error(c, http.StatusBadRequest, "invalid json format")
		return
	}

	ipAddress := c.ClientIP()
	deviceInfo := c.GetHeader("User-Agent")

	inputSignInUcase := ports.SignInUseCaseInput{
		Email:      body.Email,
		Password:   body.Password,
		IpAddress:  ipAddress,
		DeviceInfo: deviceInfo,
	}

	out, err := h.signInUseCase.SignIn(c.Request.Context(), inputSignInUcase)
	if err != nil {
		h.logger.Error("failed to sign in user", zap.Error(err))
		if errors.Is(err, usecases.ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	secureCookie := h.cfg.App.Env != "development"
	// Set cookies
	c.SetCookie("access_token", out.AccessToken, int(h.cfg.JWT.AccessTTL.Seconds()), "/", "", secureCookie, true)
	c.SetCookie("refresh_token", out.RefreshToken, int(h.cfg.JWT.RefreshTTL.Seconds()), "/", "", secureCookie, true)

	response.Success(c, http.StatusOK, gin.H{
		"message": "user signed in successfully",
	})
}
