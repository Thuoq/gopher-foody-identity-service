package ports

import (
	"context"
	"gopher-identity-service/internal/core/domain"
)

type SignInUseCaseInput struct {
	Email      string
	Password   string
	IpAddress  string
	DeviceInfo string
}

type SignInUseCaseOutput struct {
	User         *domain.User
	SessionID    string
	RefreshToken string
	AccessToken  string
}

type ISignInUseCase interface {
	SignIn(ctx context.Context, input SignInUseCaseInput) (*SignInUseCaseOutput, error)
}
