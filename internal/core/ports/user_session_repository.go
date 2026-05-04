package ports

import (
	"context"

	"gopher-identity-service/internal/core/domain"
)

type UserSessionRepository interface {
	CountActiveSessions(ctx context.Context, userID int64) (int64, error)
	DeleteOldestSession(ctx context.Context, userID int64) error
	CreateUserSession(ctx context.Context, session *domain.UserSession) error
	GetBySessionID(ctx context.Context, sessionID string) (*domain.UserSession, error)
	UpdateUserSession(ctx context.Context, session *domain.UserSession) error
	DeleteBySessionID(ctx context.Context, sessionID string) error
}
