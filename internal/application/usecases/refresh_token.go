package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"gopher-identity-service/internal/core/ports"
	"gopher-identity-service/pkg/jwt"

	"github.com/alexedwards/argon2id"
	"gorm.io/datatypes"
)

var (
	ErrTokenReuseDetected = errors.New("token reuse detected! session revoked")
	ErrInvalidToken       = errors.New("invalid refresh token")
)

type refreshTokenUseCase struct {
	userRepo        ports.UserRepository
	userSessionRepo ports.UserSessionRepository
	jwtManager      jwt.TokenManager
}

func NewRefreshTokenUseCase(
	userRepo ports.UserRepository,
	userSessionRepo ports.UserSessionRepository,
	jwtManager jwt.TokenManager,
) ports.IRefreshTokenUseCase {
	return &refreshTokenUseCase{
		userRepo:        userRepo,
		userSessionRepo: userSessionRepo,
		jwtManager:      jwtManager,
	}
}

func (uc *refreshTokenUseCase) RefreshToken(ctx context.Context, input ports.RefreshTokenUseCaseInput) (*ports.RefreshTokenUseCaseOutput, error) {
	// 1. Validate token signature and get claims
	claims, err := uc.jwtManager.ValidateRefreshToken(input.RefreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// 2. Fetch session from DB
	session, err := uc.userSessionRepo.GetBySessionID(ctx, claims.SessionID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// 3. Detect Reuse
	// Compare input token with current RefreshTokenHash
	match, err := argon2id.ComparePasswordAndHash(input.RefreshToken, session.RefreshTokenHash)
	if err != nil {
		return nil, err
	}

	if !match {
		// Potential reuse! Check history.
		var history []string
		if len(session.TokenHistory) > 0 {
			if err := json.Unmarshal(session.TokenHistory, &history); err != nil {
				return nil, err
			}
		}

		for _, oldHash := range history {
			reused, _ := argon2id.ComparePasswordAndHash(input.RefreshToken, oldHash)
			if reused {
				// REUSE DETECTED!
				_ = uc.userSessionRepo.DeleteBySessionID(ctx, claims.SessionID)
				return nil, ErrTokenReuseDetected
			}
		}

		return nil, ErrInvalidToken
	}

	// 4. Rotation
	// Generate new tokens
	accessToken, err := uc.jwtManager.GenerateAccessToken(claims.PublicUserId, claims.SessionID)
	if err != nil {
		return nil, err
	}
	refreshToken, err := uc.jwtManager.GenerateRefreshToken(claims.PublicUserId, claims.SessionID)
	if err != nil {
		return nil, err
	}

	// Hash new refresh token
	newHash, err := argon2id.CreateHash(refreshToken, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	// Update History
	var history []string
	if len(session.TokenHistory) > 0 {
		_ = json.Unmarshal(session.TokenHistory, &history)
	}
	history = append(history, session.RefreshTokenHash) // add current hash to history
	
	// Keep history size reasonable (e.g., last 10 hashes)
	if len(history) > 10 {
		history = history[len(history)-10:]
	}

	historyJSON, _ := json.Marshal(history)

	// Update Session
	session.RefreshTokenHash = newHash
	session.TokenHistory = datatypes.JSON(historyJSON)

	if err := uc.userSessionRepo.UpdateUserSession(ctx, session); err != nil {
		return nil, err
	}

	return &ports.RefreshTokenUseCaseOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
