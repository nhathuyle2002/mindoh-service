package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(r *gin.Engine, db *gorm.DB) {
	handler := NewUserHandler(db)
	r.POST("/api/register", handler.Register)
	r.POST("/api/login", handler.Login)
	r.GET("/api/users/:id", handler.GetUser)
	r.PUT("/api/users/:id", handler.UpdateUser)
	r.DELETE("/api/users/:id", handler.DeleteUser)
}
