package usecases

import (
	"context"
	"errors"
	constantUser "gopher-identity-service/internal/core/domain/constant/user"
	"gopher-identity-service/pkg/jwt"
	"time"

	"gopher-identity-service/internal/core/domain"
	"gopher-identity-service/internal/core/ports"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
)

var (
	ErrEmailExists    = errors.New("email already exists")
	ErrUsernameExists = errors.New("username already exists")
)

type signUpUseCase struct {
	userRepo   ports.UserRepository
	jwtManager jwt.TokenManager
}

func NewSignUpUseCase(
	userRepo ports.UserRepository,
	jwtManager jwt.TokenManager,
) ports.ISignUpUseCase {
	return &signUpUseCase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (uc *signUpUseCase) SignUp(ctx context.Context, input ports.SignUpUseCaseInput) (*ports.SignUpUseCaseOutput, error) {
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
		Role:     string(constantUser.RoleUser),
		PublicId: uuid.New().String(),
	}

	// 5. Generate Session and Refresh Token
	sessionID := uuid.New().String()
	accessToken, err := uc.jwtManager.GenerateAccessToken(user.PublicId, sessionID)
	if err != nil {
		return nil, err
	}
	refreshToken, err := uc.jwtManager.GenerateRefreshToken(user.PublicId, sessionID)
	if err != nil {
		return nil, err
	}
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
	return &ports.SignUpUseCaseOutput{
		User:         &user,
		SessionID:    sessionID,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}
