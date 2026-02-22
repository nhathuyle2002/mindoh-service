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

// UserLoginRequest is the request body for user login.
type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserUpdateRequest is the request body for updating user profile.
type UserUpdateRequest struct {
	Email     string `json:"email"`
	Name      string `json:"name,omitempty"`
	Birthdate string `json:"birthdate,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
}

// UserResponse is the public-facing user representation (no password hash).
type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Name      string `json:"name,omitempty"`
	Birthdate string `json:"birthdate,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
	CreatedAt string `json:"created_at"`
}

// LoginResponse is returned on successful login.
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
