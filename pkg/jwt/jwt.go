package jwt

import (
	"gopher-identity-service/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManager interface {
	GenerateAccessToken(publicUserId string, sessionId string) (string, error)
	GenerateRefreshToken(publicUserId string, sessionId string) (string, error)
}

type manager struct {
	config *config.Config
}

func NewManager(config *config.Config) TokenManager {
	return &manager{
		config: config,
	}
}

type AccessTokenClaims struct {
	PublicUserId string `json:"public_user_id"`
	SessionID    string `json:"session_id"`
	jwt.RegisteredClaims
}

func (m *manager) GenerateAccessToken(publicUserId string, sessionID string) (string, error) {
	claims := AccessTokenClaims{
		PublicUserId: publicUserId,
		SessionID:    sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.config.JWT.AccessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.JWT.AccessSecret))
}

func (m *manager) GenerateRefreshToken(publicUserId string, sessionID string) (string, error) {
	claims := AccessTokenClaims{
		PublicUserId: publicUserId,
		SessionID:    sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.config.JWT.RefreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.JWT.RefreshSecret))
}
