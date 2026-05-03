package ports

import (
	"context"
	"gopher-identity-service/internal/core/domain"
)

type SignUpUseCaseInput struct {
	Email      string
	Username   string
	Password   string
	IpAddress  string
	DeviceInfo string
}

type SignUpUseCaseOutput struct {
	User         *domain.User
	SessionID    string
	RefreshToken string
	AccessToken  string
}
type ISignUpUseCase interface {
	SignUp(ctx context.Context, x SignUpUseCaseInput) (*SignUpUseCaseOutput, error)
}
