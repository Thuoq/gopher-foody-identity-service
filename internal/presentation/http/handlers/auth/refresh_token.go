package auth

import (
	"errors"
	"net/http"

	"gopher-identity-service/internal/application/usecases"
	"gopher-identity-service/internal/core/ports"
	"gopher-identity-service/internal/presentation/http/handlers/auth/dto"
	"gopher-identity-service/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RefreshTokenHandler struct {
	refreshTokenUseCase ports.IRefreshTokenUseCase
	logger              *zap.Logger
}

func NewRefreshTokenHandler(refreshTokenUseCase ports.IRefreshTokenUseCase, logger *zap.Logger) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		refreshTokenUseCase: refreshTokenUseCase,
		logger:              logger,
	}
}

func (h *RefreshTokenHandler) Handle(c *gin.Context) {
	var body dto.RefreshTokenRequest
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

	input := ports.RefreshTokenUseCaseInput{
		RefreshToken: body.RefreshToken,
	}

	out, err := h.refreshTokenUseCase.RefreshToken(c.Request.Context(), input)
	if err != nil {
		h.logger.Error("failed to refresh token", zap.Error(err))
		
		if errors.Is(err, usecases.ErrTokenReuseDetected) {
			response.Error(c, http.StatusForbidden, err.Error())
			return
		}
		
		if errors.Is(err, usecases.ErrInvalidToken) {
			response.Error(c, http.StatusUnauthorized, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	// Returns tokens directly in JSON as requested
	response.Success(c, http.StatusOK, gin.H{
		"access_token":  out.AccessToken,
		"refresh_token": out.RefreshToken,
	})
}
