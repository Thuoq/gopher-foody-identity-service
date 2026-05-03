package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"

	"gopher-identity-service/internal/core/domain"
	"gopher-identity-service/internal/core/ports"
)

var (
	ErrEmailExists    = errors.New("email already exists")
	ErrUsernameExists = errors.New("username already exists")
)

type SignUpUseCase struct {
	userRepo ports.UserRepository
}

func NewSignUpUseCase(
	userRepo ports.UserRepository,
) SignUpUseCase {
	return SignUpUseCase{
		userRepo: userRepo,
	}
}

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
}

func (uc *SignUpUseCase) SignUp(ctx context.Context, input SignUpUseCaseInput) (*SignUpUseCaseOutput, error) {
	// 1. Check if email exists
	emailExists, err := uc.userRepo.CheckEmailExists(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, ErrEmailExists
	}

	// 2. Check if username exists
	usernameExists, err := uc.userRepo.CheckUsernameExists(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, ErrUsernameExists
	}

	// 3. Hash password using Argon2id
	hashedPassword, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	// 4. Create User domain model
	user := domain.User{
		Email:    input.Email,
		Username: input.Username,
		Password: hashedPassword,
		Role:     "user",
	}

	// 5. Generate Session and Refresh Token
	sessionID := uuid.New().String()
	refreshToken := uuid.New().String() + "-" + uuid.New().String()

	// Hash refresh token for storage
	hashedRefreshToken, err := argon2id.CreateHash(refreshToken, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	// 6. Create UserSession domain model
	session := domain.UserSession{
		SessionId:        sessionID,
		RefreshTokenHash: hashedRefreshToken,
		IpAddress:        input.IpAddress,
		DeviceInfo:       input.DeviceInfo,
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour), // e.g. 7 days
	}

	// 7. Save both in a transaction
	if err := uc.userRepo.CreateUserWithSession(ctx, &user, &session); err != nil {
		return nil, err
	}

	// 8. Return success
	return &SignUpUseCaseOutput{
		User:         &user,
		SessionID:    sessionID,
		RefreshToken: refreshToken,
	}, nil
}
