package repositories

import (
	"context"

	"gopher-identity-service/internal/core/domain"
	"gopher-identity-service/internal/core/ports"

	"gorm.io/gorm"
)

type userSessionRepository struct {
	db *gorm.DB
}

func NewUserSessionPostgresRepository(db *gorm.DB) ports.UserSessionRepository {
	return &userSessionRepository{
		db: db,
	}
}

func (r *userSessionRepository) CountActiveSessions(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.UserSession{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *userSessionRepository) DeleteOldestSession(ctx context.Context, userID int64) error {
	// Find the oldest session ID for the user
	var oldestSession domain.UserSession
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at asc").
		First(&oldestSession).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // No sessions to delete
		}
		return err
	}

	// Delete it
	return r.db.WithContext(ctx).Delete(&oldestSession).Error
}

func (r *userSessionRepository) CreateUserSession(ctx context.Context, session *domain.UserSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *userSessionRepository) GetBySessionID(ctx context.Context, sessionID string) (*domain.UserSession, error) {
	var session domain.UserSession
	err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *userSessionRepository) UpdateUserSession(ctx context.Context, session *domain.UserSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

func (r *userSessionRepository) DeleteBySessionID(ctx context.Context, sessionID string) error {
	return r.db.WithContext(ctx).Where("session_id = ?", sessionID).Delete(&domain.UserSession{}).Error
}

