package user

import (
	dbmodel "mindoh-service/internal/db"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *dbmodel.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) GetByID(id uint) (*dbmodel.User, error) {
	var user dbmodel.User
	err := r.DB.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) GetByUsername(username string) (*dbmodel.User, error) {
	var user dbmodel.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) Update(user *dbmodel.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) UpdateFields(userID uint, fields map[string]interface{}) error {
	return r.DB.Model(&dbmodel.User{}).Where("id = ?", userID).Updates(fields).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.DB.Delete(&dbmodel.User{}, id).Error
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
