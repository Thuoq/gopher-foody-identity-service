package usecases

import (
	"context"
	"errors"
	"gopher-identity-service/internal/config"
	"gopher-identity-service/internal/core/domain"
	"gopher-identity-service/internal/core/ports"
	"gopher-identity-service/pkg/jwt"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type signInUseCase struct {
	userRepo        ports.UserRepository
	userSessionRepo ports.UserSessionRepository
	jwtManager      jwt.TokenManager
	cfg             *config.Config
}

func NewSignInUseCase(
	userRepo ports.UserRepository,
	userSessionRepo ports.UserSessionRepository,
	jwtManager jwt.TokenManager,
	cfg *config.Config,
) ports.ISignInUseCase {
	return &signInUseCase{
		userRepo:        userRepo,
		userSessionRepo: userSessionRepo,
		jwtManager:      jwtManager,
		cfg:             cfg,
	}
}

func (uc *signInUseCase) SignIn(ctx context.Context, input ports.SignInUseCaseInput) (*ports.SignInUseCaseOutput, error) {
	// 1. Fetch user by email
	user, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 2. Verify password
	match, err := argon2id.ComparePasswordAndHash(input.Password, user.Password)
	if err != nil || !match {
		return nil, ErrInvalidCredentials
	}

	// 3. Enforce 3-device limit
	count, err := uc.userSessionRepo.CountActiveSessions(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if count >= 3 {
		err = uc.userSessionRepo.DeleteOldestSession(ctx, user.ID)
		if err != nil {
			return nil, err
		}
	}

	// 4. Generate tokens
	sessionID := uuid.New().String()
	accessToken, err := uc.jwtManager.GenerateAccessToken(user.PublicId, sessionID)
	if err != nil {
		return nil, err
	}
	refreshToken, err := uc.jwtManager.GenerateRefreshToken(user.PublicId, sessionID)
	if err != nil {
		return nil, err
	}

	// 5. Hash refresh token
	hashedRefreshToken, err := argon2id.CreateHash(refreshToken, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	// 6. Create UserSession model
	session := domain.UserSession{
		UserId:           user.ID,
		SessionId:        sessionID,
		RefreshTokenHash: hashedRefreshToken,
		IpAddress:        input.IpAddress,
		DeviceInfo:       input.DeviceInfo,
		ExpiresAt:        time.Now().Add(uc.cfg.JWT.RefreshTTL),
	}

	// 7. Save UserSession
	err = uc.userSessionRepo.CreateUserSession(ctx, &session)
	if err != nil {
		return nil, err
	}

	// 8. Return success
	return &ports.SignInUseCaseOutput{
		User:         user,
		SessionID:    sessionID,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}
