package main

import (
	"mindoh-service/config"
	"mindoh-service/internal/auth"
	"mindoh-service/internal/currency"
	"mindoh-service/internal/db"
	"mindoh-service/internal/expense"
	"mindoh-service/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	_ "mindoh-service/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Mindoh Service API
// @version         1.0
// @description     A personal finance management API service
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// Services holds all the service instances for the application
type Services struct {
	Config         *config.Config
	DB             *gorm.DB
	UserService    *user.UserService
	AuthService    auth.IAuthService
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
	authService := auth.NewAuthService(cfg)

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
	user.RegisterUserRoutes(r, s.AuthService, s.UserService)
	// Register expense routes
	expense.RegisterExpenseRoutes(r, s.AuthService, s.ExpenseService)
	// Register currency routes
	currency.RegisterCurrencyRoutes(r, s.AuthService)
}

func main() {
	// Initialize all services
	services := NewService()

	r := gin.Default()
	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register application routes
	RegisterRoutes(r, services)

	r.Run()
}
