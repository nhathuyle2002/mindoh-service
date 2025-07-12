package user

import "gorm.io/gorm"

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"unique;not null" json:"username"`
	Email        string         `gorm:"unique;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Role         string         `gorm:"default:user" json:"role"` // user or admin
	CreatedAt    int64          `json:"created_at"`
	UpdatedAt    int64          `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
