package user

import (
	"errors"
	"net/http"

	"gopher-identity-service/internal/application/usecases"
	"gopher-identity-service/internal/presentation/http/handlers/user/dto/request"
	"gopher-identity-service/pkg/jwt"
	"gopher-identity-service/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SignUpHandler struct {
	signUpUseCase usecases.SignUpUseCase
	jwtManager    jwt.TokenManager
	logger        *zap.Logger
}

func NewSignUpHandler(signUpUseCase usecases.SignUpUseCase, jwtManager jwt.TokenManager, logger *zap.Logger) SignUpHandler {
	return SignUpHandler{
		signUpUseCase: signUpUseCase,
		jwtManager:    jwtManager,
		logger:        logger,
	}
}

func (h *SignUpHandler) Handle(c *gin.Context) {
	var body request.SignUpRequest
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

	inputSignUpUcase := usecases.SignUpUseCaseInput{
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

	accessToken, err := h.jwtManager.GenerateAccessToken(out.User.ID, out.SessionID)
	if err != nil {
		h.logger.Error("failed to generate access token", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	// Set cookies
	// httpOnly, secure, domain, path, maxAge
	c.SetCookie("access_token", accessToken, int(15*60), "/", "", false, true) // 15 mins
	c.SetCookie("refresh_token", out.RefreshToken, int(7*24*60*60), "/", "", false, true) // 7 days

	response.Success(c, http.StatusCreated, gin.H{
		"message": "user created successfully",
	})
}

