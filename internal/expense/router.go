package expense

import (
	"mindoh-service/internal/auth"

	"github.com/gin-gonic/gin"
)

func RegisterExpenseRoutes(r *gin.Engine, a auth.IAuthService, service *ExpenseService) {
	handler := NewExpenseHandler(service)

	group := r.Group("/api/expenses")
	group.Use(a.AuthMiddleware())
	{
		group.POST("/", handler.AddExpense)
		group.PUT("/:id", handler.UpdateExpense)
		group.GET("/", handler.ListExpenses)
		group.GET("/summary", handler.Summary)
		group.GET("/exchange-rates", handler.GetExchangeRates)
	}
}
