package auth

import (
	"errors"
	"gopher-identity-service/internal/config"
	"gopher-identity-service/internal/core/ports"
	"gopher-identity-service/internal/presentation/http/handlers/auth/dto"
	"net/http"

	"gopher-identity-service/internal/application/usecases"
	"gopher-identity-service/pkg/jwt"
	"gopher-identity-service/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SignUpHandler struct {
	signUpUseCase ports.ISignUpUseCase
	jwtManager    jwt.TokenManager
	cfg           *config.Config
	logger        *zap.Logger
}

func NewSignUpHandler(signUpUseCase ports.ISignUpUseCase, jwtManager jwt.TokenManager, logger *zap.Logger) *SignUpHandler {
	return &SignUpHandler{
		signUpUseCase: signUpUseCase,
		jwtManager:    jwtManager,
		logger:        logger,
	}
}

func (h *SignUpHandler) Handle(c *gin.Context) {
	var body dto.SignUpRequest
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

	inputSignUpUcase := ports.SignUpUseCaseInput{
		Email:      body.Email,
		Password:   body.Password,
		Username:   body.Username,
		IpAddress:  ipAddress,
		DeviceInfo: deviceInfo,
	}

	out, err := h.signUpUseCase.SignUp(c.Request.Context(), inputSignUpUcase)
	if err != nil {
		h.logger.Error("failed to sign up user", zap.Error(err))
		if errors.Is(err, usecases.ErrEmailExists) || errors.Is(err, usecases.ErrUsernameExists) {
			response.Error(c, http.StatusConflict, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	secureCookie := h.cfg.App.Env != "development"
	// Set cookies
	// httpOnly, secure, domain, path, maxAge
	c.SetCookie("access_token", out.AccessToken, int(h.cfg.JWT.AccessTTL.Seconds()), "/", "", secureCookie, true)
	c.SetCookie("refresh_token", out.RefreshToken, int(h.cfg.JWT.RefreshTTL.Seconds()), "/", "", secureCookie, true)

	response.Success(c, http.StatusCreated, gin.H{
		"message": "user created successfully",
	})
}
