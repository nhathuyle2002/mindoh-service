package expense

import (
	"net/http"
	"time"

	"mindoh-service/common/utils"
	"mindoh-service/internal/auth"

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
	authCtx := auth.GetAuthContext(c)
	userID := authCtx.UserID
	role := authCtx.Role
	if role == auth.RoleUser && req.UserID != 0 && req.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only add your own expenses"})
		return
	}
	if role == auth.RoleUser || req.UserID == 0 {
		req.UserID = userID
	}
	// Set date to now if not provided
	if req.Date.IsZero() {
		req.Date = time.Now()
	}
	// Normalize date to start of day (00:00:00)
	req.Date = time.Date(req.Date.Year(), req.Date.Month(), req.Date.Day(), 0, 0, 0, 0, req.Date.Location())
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

// UpdateExpense godoc
// @Summary Update an existing expense
// @Description Update details of an existing expense
// @Tags expenses
// @Accept json
// @Produce json
// @Param id path int true "Expense ID"
// @Param expense body ExpenseUpdateRequest true "Expense update details"
// @Success 200 {object} Expense "Expense updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Expense not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /expenses/{id} [put]
func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	var req ExpenseUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	authCtx := auth.GetAuthContext(c)
	userID := authCtx.UserID
	expenseID := utils.ParseUint(c.Param("id"))

	// Fetch existing expense to check ownership
	expense, err := h.Service.Repo.GetByID(expenseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	role := authCtx.Role
	if role == auth.RoleUser && expense.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only add your own expenses"})
		return
	}

	// Update fields if provided
	if req.Amount != nil {
		expense.Amount = *req.Amount
	}
	if req.Currency != nil {
		expense.Currency = *req.Currency
	}
	if req.Kind != nil {
		expense.Kind = *req.Kind
	}
	if req.Type != nil {
		expense.Type = *req.Type
	}
	if req.Description != nil {
		expense.Description = *req.Description
	}
	if req.Date != nil {
		// Normalize date to start of day (00:00:00)
		normalizedDate := time.Date(req.Date.Year(), req.Date.Month(), req.Date.Day(), 0, 0, 0, 0, req.Date.Location())
		expense.Date = normalizedDate
	}

	if err := h.Service.UpdateExpense(expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense"})
		return
	}
	c.JSON(http.StatusOK, expense)
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
	authCtx := auth.GetAuthContext(c)
	userID := authCtx.UserID
	role := authCtx.Role
	var filter ExpenseFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}
	if role == auth.RoleUser && filter.UserID != 0 && filter.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own expenses"})
		return
	}
	if role == auth.RoleUser || filter.UserID == 0 {
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
	authCtx := auth.GetAuthContext(c)
	userID := authCtx.UserID
	role := authCtx.Role
	var filter ExpenseFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}
	if role == auth.RoleUser && filter.UserID != 0 && filter.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own expenses"})
		return
	}
	if role == auth.RoleUser || filter.UserID == 0 {
		filter.UserID = userID
	}

	summary, err := h.Service.Summary(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch summary"})
		return
	}
	c.JSON(http.StatusOK, summary)
}

func (h *ExpenseHandler) GetExchangeRates(c *gin.Context) {
	rates := h.Service.GetExchangeRates()
	c.JSON(http.StatusOK, gin.H{
		"base_currency": "VND",
		"rates":         rates,
	})
}

// GetAvailableCurrencies godoc
// @Summary Get available currencies
// @Description Get list of available currencies
// @Tags expenses
// @Produce json
// @Success 200 {object} map[string]interface{} "List of available currencies"
// @Security BearerAuth
// @Router /expenses/currencies [get]
func (h *ExpenseHandler) GetAvailableCurrencies(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"currencies": AvailableCurrencies,
	})
}
