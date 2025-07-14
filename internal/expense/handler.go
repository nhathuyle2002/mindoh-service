package expense

import (
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

// AddExpense godoc
// @Summary Add a new expense
// @Description Create a new expense record
// @Tags expenses
// @Accept json
// @Produce json
// @Param expense body ExpenseCreateRequest true "Expense details"
// @Success 201 {object} Expense "Expense created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /expenses [post]
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

// ListExpenses godoc
// @Summary List expenses
// @Description Get list of expenses with optional filtering
// @Tags expenses
// @Accept json
// @Produce json
// @Param user_id query int false "User ID"
// @Param kind query string false "Expense kind (expense/income)"
// @Param type query string false "Expense type (food/salary/transport/entertainment)"
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} Expense "List of expenses"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /expenses [get]
func (h *ExpenseHandler) ListExpenses(c *gin.Context) {
	authCtx := user.GetAuthContext(c)
	userID := authCtx.UserID
	role := authCtx.Role
	var filter ExpenseFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}
	if role == user.RoleUser && filter.UserID != 0 && filter.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own expenses"})
		return
	}
	if role == user.RoleUser || filter.UserID == 0 {
		filter.UserID = userID
	}

	expenses, err := h.Service.ListExpenses(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}
	c.JSON(http.StatusOK, expenses)
}

// Summary godoc
// @Summary Get expense summary
// @Description Get summary of expenses with optional filtering
// @Tags expenses
// @Accept json
// @Produce json
// @Param user_id query int false "User ID"
// @Param kind query string false "Expense kind (expense/income)"
// @Param type query string false "Expense type (food/salary/transport/entertainment)"
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} ExpenseSummary "Expense summary"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /expenses/summary [get]
func (h *ExpenseHandler) Summary(c *gin.Context) {
	authCtx := user.GetAuthContext(c)
	userID := authCtx.UserID
	role := authCtx.Role
	var filter ExpenseFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}
	if role == user.RoleUser && filter.UserID != 0 && filter.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own expenses"})
		return
	}
	if role == user.RoleUser || filter.UserID == 0 {
		filter.UserID = userID
	}

	summary, err := h.Service.Summary(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch summary"})
		return
	}
	c.JSON(http.StatusOK, summary)
}
