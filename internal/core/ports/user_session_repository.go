package ports

import (
	"context"

	"gopher-identity-service/internal/core/domain"
)

type UserSessionRepository interface {
	CountActiveSessions(ctx context.Context, userID int64) (int64, error)
	DeleteOldestSession(ctx context.Context, userID int64) error
	CreateUserSession(ctx context.Context, session *domain.UserSession) error
}
