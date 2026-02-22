package user

import (
	"mindoh-service/common/utils"
	"mindoh-service/internal/auth"
	dbmodel "mindoh-service/internal/db"
	"mindoh-service/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	authService auth.IAuthService
	userService *UserService
}

func NewUserHandler(authService auth.IAuthService, userService *UserService) *UserHandler {
	return &UserHandler{
		authService: authService,
		userService: userService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.UserRegisterRequest true "User registration details"
// @Success 201 {object} dto.UserResponse "User created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user := &dbmodel.User{
		Username:  req.Username,
		Email:     req.Email,
		Name:      req.Name,
		Birthdate: req.Birthdate,
		Phone:     req.Phone,
		Address:   req.Address,
		Role:      auth.RoleUser,
	}
	user.PasswordHash, _ = HashPassword(req.Password)
	if err := h.userService.CreateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists or invalid data"})
		return
	}
	c.JSON(http.StatusCreated, toUserResponse(user))
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.UserLoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user, err := h.userService.ValidateCredentials(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token, err := h.authService.GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, dto.LoginResponse{Token: token, User: toUserResponse(user)})
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get user information by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse "User found"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetUserByID(utils.ParseUint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(user))
}

// UpdateUser godoc
// @Summary Update user information
// @Description Update user information by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body dto.UserUpdateRequest true "User update details"
// @Success 200 {object} dto.UserResponse "User updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Security BearerAuth
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user, err := h.userService.GetUserByID(utils.ParseUint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	// Update fields
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Birthdate != "" {
		user.Birthdate = req.Birthdate
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Address != "" {
		user.Address = req.Address
	}
	if err := h.userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(user))
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{} "User deleted successfully"
// @Failure 500 {object} map[string]interface{} "Failed to delete user"
// @Security BearerAuth
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.userService.DeleteUser(utils.ParseUint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// AdminCreateUser godoc
// @Summary Create user (admin)
// @Description Admin creates a user with a specified role
// @Tags admin
// @Accept json
// @Produce json
// @Param user body dto.AdminCreateUserRequest true "User details"
// @Success 201 {object} dto.UserResponse "User created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Security BearerAuth
// @Router /admin/users [post]
func (h *UserHandler) AdminCreateUser(c *gin.Context) {
	var req dto.AdminCreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role := auth.Role(req.Role)
	if role == "" {
		role = auth.RoleUser
	}
	user := &dbmodel.User{
		Username:  req.Username,
		Email:     req.Email,
		Name:      req.Name,
		Birthdate: req.Birthdate,
		Phone:     req.Phone,
		Address:   req.Address,
		Role:      role,
	}
	user.PasswordHash, _ = HashPassword(req.Password)
	if err := h.userService.CreateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists or invalid data"})
		return
	}
	c.JSON(http.StatusCreated, toUserResponse(user))
}
