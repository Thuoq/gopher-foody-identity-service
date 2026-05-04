package usecases

import (
	"context"
	"gopher-identity-service/internal/core/ports"
)

type logoutUseCase struct {
	userSessionRepo ports.UserSessionRepository
}

func NewLogoutUseCase(userSessionRepo ports.UserSessionRepository) ports.ILogoutUseCase {
	return &logoutUseCase{
		userSessionRepo: userSessionRepo,
	}
}

func (uc *logoutUseCase) Logout(ctx context.Context, input ports.LogoutUseCaseInput) error {
	return uc.userSessionRepo.DeleteBySessionID(ctx, input.SessionID)
}
