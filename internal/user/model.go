package user

import (
	"mindoh-service/internal/auth"
	"time"

	"gorm.io/gorm"
)

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
}

type UserRegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Name      string `json:"name,omitempty"`
	Birthdate string `json:"birthdate,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	Email     string `json:"email"`
	Name      string `json:"name,omitempty"`
	Birthdate string `json:"birthdate,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
}
