package user

import (
	"fmt"

	"gorm.io/gorm"
)

// UserService handles business logic for users
type UserService struct {
	DB   *gorm.DB
	Repo *UserRepository
}

// NewUserService creates a new user service
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		DB:   db,
		Repo: NewUserRepository(db),
	}
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(id uint) (*User, error) {
	return s.Repo.GetByID(id)
}

// CreateUser creates a new user
func (s *UserService) CreateUser(user *User) error {
	return s.Repo.Create(user)
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(user *User) error {
	return s.Repo.Update(user)
}

// DeleteUser deletes a user by their ID
func (s *UserService) DeleteUser(id uint) error {
	return s.Repo.Delete(id)
}

// ValidateCredentials validates user login credentials
func (s *UserService) ValidateCredentials(username, password string) (*User, error) {
	user, err := s.Repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}
