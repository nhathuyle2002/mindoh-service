package expense

import (
	"fmt"
	"net/http"
	"time"

	"mindoh-service/internal/user"

	"github.com/gin-gonic/gin"
)

// ExpenseHandler handles HTTP requests for expenses
type ExpenseHandler struct {
	Service *ExpenseService
}

func NewExpenseHandler(service *ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{Service: service}
}

func (h *ExpenseHandler) AddExpense(c *gin.Context) {
	var req ExpenseCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	authCtx := user.GetAuthContext(c)
	userID := authCtx.UserID
	role := authCtx.Role
	if role == user.RoleUser && req.UserID != 0 && req.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only add your own expenses"})
		return
	}
	if role == user.RoleUser || req.UserID == 0 {
		req.UserID = userID
	}
	// Set date to now if not provided
	if req.Date.IsZero() {
		req.Date = time.Now()
	}
	expense := Expense{
		UserID:      req.UserID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Kind:        req.Kind,
		Type:        req.Type,
		Description: req.Description,
		Date:        req.Date,
	}
	if err := h.Service.AddExpense(&expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add expense"})
		return
	}
	c.JSON(http.StatusCreated, expense)
}

func (h *ExpenseHandler) ListExpenses(c *gin.Context) {
	authCtx := user.GetAuthContext(c)
	userID := authCtx.UserID
	role := authCtx.Role
	kind := c.Query("kind")
	typeStr := c.Query("type")
	var expenses []Expense
	var err error
	if role == user.RoleUser {
		expenses, err = h.Service.ListExpenses(userID, kind, typeStr)
	} else {
		// admin can list all users' expenses if user_id is provided as query param
		userIDParam := c.Query("user_id")
		if userIDParam != "" {
			// parse userIDParam to uint
			var uid uint
			_, errParse := fmt.Sscanf(userIDParam, "%d", &uid)
			if errParse == nil {
				expenses, err = h.Service.ListExpenses(uid, kind, typeStr)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id param"})
				return
			}
		} else {
			// list all expenses for all users
			expenses, err = h.Service.ListExpenses(0, kind, typeStr)
		}
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}
	c.JSON(http.StatusOK, expenses)
}

// Example: GET /expenses/summary/day?date=2025-07-13&kind=expense&type=food
func (h *ExpenseHandler) SummaryByDay(c *gin.Context) {
	authCtx := user.GetAuthContext(c)
	userID := authCtx.UserID
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing date parameter"})
		return
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format (use YYYY-MM-DD)"})
		return
	}
	kind := c.Query("kind")
	typeStr := c.Query("type")
	sum, err := h.Service.SummaryByDay(userID, date, kind, typeStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to summarize expenses"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": sum})
}
