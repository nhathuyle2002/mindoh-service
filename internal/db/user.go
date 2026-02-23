package db

import (
	"mindoh-service/internal/auth"
	"time"

	"gorm.io/gorm"
)

// User is the database model for a user account.
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"unique;not null" json:"username"`
	Email        string         `gorm:"unique;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Role         auth.Role      `gorm:"default:user" json:"role"`
	Name         string         `json:"name,omitempty"`
	Birthdate    string         `json:"birthdate,omitempty"`
	Phone        string         `json:"phone,omitempty"`
	Address      string         `json:"address,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Email verification
	IsEmailVerified   bool      `gorm:"default:false" json:"is_email_verified"`
	EmailVerifyToken  string    `gorm:"index" json:"-"`
	EmailVerifyExpiry time.Time `json:"-"`

	// Password reset
	PasswordResetToken  string    `gorm:"index" json:"-"`
	PasswordResetExpiry time.Time `json:"-"`
}
