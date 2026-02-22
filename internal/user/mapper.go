package user

import (
	dbmodel "mindoh-service/internal/db"
	"mindoh-service/internal/dto"
)

// toUserResponse maps a User model to a UserResponse DTO.
func toUserResponse(u *dbmodel.User) dto.UserResponse {
	return dto.UserResponse{
		Username:  u.Username,
		Email:     u.Email,
		Role:      string(u.Role),
		Name:      u.Name,
		Birthdate: u.Birthdate,
		Phone:     u.Phone,
		Address:   u.Address,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
