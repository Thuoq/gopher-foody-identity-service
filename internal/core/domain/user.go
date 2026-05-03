package domain

import (
	"time"
)

type User struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	PublicId  string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Username  string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password  string    `gorm:"type:varchar(255);not null"`
	Role      string    `json:"role" gorm:"type:varchar(255);default:'user'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
