package ports

import "context"

type LogoutUseCaseInput struct {
	SessionID string
}

type ILogoutUseCase interface {
	Logout(ctx context.Context, input LogoutUseCaseInput) error
}
