package expense

import (
	"mindoh-service/config"
	"mindoh-service/internal/user"

	"github.com/gin-gonic/gin"
)

func RegisterExpenseRoutes(r *gin.Engine, cfg *config.Config, service *ExpenseService) {
	handler := NewExpenseHandler(service)

	group := r.Group("/api/expenses")
	group.Use(user.AuthMiddleware(cfg))
	{
		group.POST("/", handler.AddExpense)
		group.GET("/", handler.ListExpenses)
		group.GET("/summary", handler.Summary)
	}
}
