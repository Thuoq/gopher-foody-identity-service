package ports

import (
	"context"

	"gopher-identity-service/internal/core/domain"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	// New methods for Sign Up
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	CreateUserWithSession(ctx context.Context, user *domain.User, session *domain.UserSession) error
}

