package auth

import (
	"gopher-identity-service/internal/core/ports"
	"gopher-identity-service/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LogoutHandler struct {
	logoutUseCase ports.ILogoutUseCase
	logger        *zap.Logger
}

func NewLogoutHandler(logoutUseCase ports.ILogoutUseCase, logger *zap.Logger) *LogoutHandler {
	return &LogoutHandler{
		logoutUseCase: logoutUseCase,
		logger:        logger,
	}
}

func (h *LogoutHandler) Handle(c *gin.Context) {
	sessionID, exists := c.Get("session_id")
	if !exists {
		response.Error(c, http.StatusInternalServerError, "session not found in context")
		return
	}

	input := ports.LogoutUseCaseInput{
		SessionID: sessionID.(string),
	}

	if err := h.logoutUseCase.Logout(c.Request.Context(), input); err != nil {
		h.logger.Error("failed to logout", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, "failed to revoke session")
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"message": "logged out successfully",
	})
}
