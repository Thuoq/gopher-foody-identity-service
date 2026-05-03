package repositories

import (
	"context"
	"gopher-identity-service/internal/core/domain"
	"gopher-identity-service/internal/core/ports"

	"gorm.io/gorm"
)

type userPostgresRepo struct {
	db *gorm.DB
}

func NewUserPostgresRepository(db *gorm.DB) ports.UserRepository {
	return &userPostgresRepo{
		db: db,
	}
}

func (r *userPostgresRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userPostgresRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userPostgresRepo) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *userPostgresRepo) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (r *userPostgresRepo) CreateUserWithSession(ctx context.Context, user *domain.User, session *domain.UserSession) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Create User
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// 2. Assign the generated User ID to Session
		session.UserId = user.ID

		// 3. Create Session
		if err := tx.Create(session).Error; err != nil {
			return err
		}

		return nil
	})
}
