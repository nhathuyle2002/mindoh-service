package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	authService IAuthService
	Repo        *UserRepository
}

func NewUserHandler(authService IAuthService, db *gorm.DB) *UserHandler {
	return &UserHandler{
		authService: authService,
		Repo:        NewUserRepository(db),
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	passwordHash, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user := &User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Name:         req.Name,
		Birthdate:    req.Birthdate,
		Phone:        req.Phone,
		Address:      req.Address,
		Role:         RoleUser, // Default role is user
	}
	if err := h.Repo.Create(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists or invalid data"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user, err := h.Repo.GetByUsername(req.Username)
	if err != nil || !CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token, err := h.authService.GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := h.Repo.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := h.Repo.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	var req UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	fields := make(map[string]interface{})
	if req.Email != "" {
		fields["email"] = req.Email
	}
	if req.Name != "" {
		fields["name"] = req.Name
	}
	if req.Birthdate != "" {
		fields["birthdate"] = req.Birthdate
	}
	if req.Phone != "" {
		fields["phone"] = req.Phone
	}
	if req.Address != "" {
		fields["address"] = req.Address
	}

	if len(fields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	if err := h.Repo.UpdateFields(user.ID, fields); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Fetch updated user
	h.Repo.DB.First(&user, id)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.Repo.Delete(parseUint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func parseUint(s string) uint {
	var i uint
	fmt.Sscanf(s, "%d", &i)
	return i
}
