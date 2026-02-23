package dto

// UserRegisterRequest is the request body for registering a new user.
type UserRegisterRequest struct {
	Username  string `json:"username"  binding:"required"`
	Email     string `json:"email"     binding:"required,email"`
	Password  string `json:"password"  binding:"required,min=6"`
	Name      string `json:"name,omitempty"`
	Birthdate string `json:"birthdate,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
}

// AdminCreateUserRequest is the request body for admin creating a user with a specific role.
type AdminCreateUserRequest struct {
	Username  string `json:"username"  binding:"required"`
	Email     string `json:"email"     binding:"required,email"`
	Password  string `json:"password"  binding:"required,min=6"`
	Role      string `json:"role"      binding:"omitempty,oneof=admin user"` // defaults to "user" if omitted
	Name      string `json:"name,omitempty"`
	Birthdate string `json:"birthdate,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
}

// UserLoginRequest is the request body for user login.
type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserUpdateRequest is the request body for updating user profile.
type UserUpdateRequest struct {
	Name      string `json:"name,omitempty"`
	Birthdate string `json:"birthdate,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
}

// UpdateEmailRequest is the request body for changing a user's email.
type UpdateEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// UserResponse is the public-facing user representation (no password hash).
type UserResponse struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Role            string `json:"role"`
	IsEmailVerified bool   `json:"is_email_verified"`
	Name            string `json:"name,omitempty"`
	Birthdate       string `json:"birthdate,omitempty"`
	Phone           string `json:"phone,omitempty"`
	Address         string `json:"address,omitempty"`
	CreatedAt       string `json:"created_at"`
}

// LoginResponse is returned on successful login.
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// ForgotPasswordRequest is the request body for initiating a password reset.
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest is the request body for completing a password reset.
type ResetPasswordRequest struct {
	Token    string `json:"token"    binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// ChangePasswordRequest is the request body for changing a password while authenticated.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password"     binding:"required,min=6"`
}
