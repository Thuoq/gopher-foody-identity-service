package ports

import "context"

type RefreshTokenUseCaseInput struct {
	RefreshToken string
}

type RefreshTokenUseCaseOutput struct {
	AccessToken  string
	RefreshToken string
}

type IRefreshTokenUseCase interface {
	RefreshToken(ctx context.Context, input RefreshTokenUseCaseInput) (*RefreshTokenUseCaseOutput, error)
}
