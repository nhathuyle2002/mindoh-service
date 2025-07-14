package main

import (
	"mindoh-service/config"
	"mindoh-service/internal/db"
	"mindoh-service/internal/expense"
	"mindoh-service/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// Services holds all the service instances for the application
type Services struct {
	Config         *config.Config
	DB             *gorm.DB
	UserService    *user.UserService
	AuthService    user.IAuthService
	ExpenseService *expense.ExpenseService
}

// NewService initializes all services for the application
func NewService() *Services {
	// Load environment variables
	godotenv.Load()

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db.ConnectDatabase(cfg)
	dbInstance := db.GetDB()

	// Initialize auth service
	authService := user.NewAuthService(cfg)

	// Initialize user service
	userService := user.NewUserService(dbInstance)

	// Initialize expense service
	expenseRepo := expense.NewExpenseRepository(dbInstance)
	expenseService := expense.NewExpenseService(expenseRepo)

	return &Services{
		Config:         cfg,
		DB:             dbInstance,
		AuthService:    authService,
		UserService:    userService,
		ExpenseService: expenseService,
	}
}

func RegisterRoutes(r *gin.Engine, s *Services) {
	// Register user routes
	user.RegisterUserRoutes(r, s.Config, s.AuthService, s.UserService)
	// Register expense routes
	expense.RegisterExpenseRoutes(r, s.Config, s.ExpenseService)
}

func main() {
	// Initialize all services
	services := NewService()

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Register application routes
	RegisterRoutes(r, services)

	r.Run()
}
