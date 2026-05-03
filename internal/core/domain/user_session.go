package domain

import (
	"time"

	"gorm.io/datatypes"
)

type UserSession struct {
	ID               int64          `gorm:"primaryKey"`
	UserId           int64          `gorm:"not null; index"`
	User             User           `gorm:"foreignKey:UserId;references:ID;constraint:OnDelete:CASCADE"`
	SessionId        string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	RefreshTokenHash string         `gorm:"type:varchar(255);not null"`
	TokenHistory     datatypes.JSON `gorm:"type:jsonb"`
	DeviceInfo       string         `gorm:"type:varchar(255)"`
	IpAddress        string         `gorm:"type:varchar(255)"`
	ExpiresAt        time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
