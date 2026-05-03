package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManager interface {
	GenerateAccessToken(userID int64, sessionID string) (string, error)
	GenerateRefreshToken() string
}

type manager struct {
	secretKey     string
	accessTTL     time.Duration
}

func NewManager(secretKey string, accessTTL time.Duration) TokenManager {
	return &manager{
		secretKey: secretKey,
		accessTTL: accessTTL,
	}
}

type AccessTokenClaims struct {
	UserID    int64  `json:"user_id"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

func (m *manager) GenerateAccessToken(userID int64, sessionID string) (string, error) {
	claims := AccessTokenClaims{
		UserID:    userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *manager) GenerateRefreshToken() string {
	// A simple random secure string is usually enough for refresh token, using UUID for simplicity and uniqueness.
	// You can also use crypto/rand for more entropy if needed.
	return uuid.New().String() + "-" + uuid.New().String()
}
