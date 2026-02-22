package expense

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"mindoh-service/common/utils"
	"mindoh-service/internal/auth"
	dbmodel "mindoh-service/internal/db"
	"mindoh-service/internal/dto"

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
// @Param expense body dto.ExpenseCreateRequest true "Expense details"
// @Success 201 {object} Expense "Expense created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /expenses [post]
func (h *ExpenseHandler) AddExpense(c *gin.Context) {
	var req dto.ExpenseCreateRequest
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
	// Set date to today if not provided, format: YYYY-MM-DD
	if req.Date == "" {
		req.Date = time.Now().Format("2006-01-02")
	}
	// Validate date format
	if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, expected YYYY-MM-DD"})
		return
	}
	expense := dbmodel.Expense{
		UserID:      req.UserID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Kind:        dbmodel.ExpenseKind(req.Kind),
		Type:        strings.ToLower(strings.TrimSpace(req.Type)),
		Resource:    dbmodel.ExpenseResource(req.Resource),
		Description: req.Description,
		Date:        req.Date,
	}
	if err := h.Service.AddExpense(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, toExpenseResponse(&expense))
}

// UpdateExpense godoc
// @Summary Update an existing expense
// @Description Update details of an existing expense
// @Tags expenses
// @Accept json
// @Produce json
// @Param id path int true "Expense ID"
// @Param expense body dto.ExpenseUpdateRequest true "Expense update details"
// @Success 200 {object} Expense "Expense updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Expense not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /expenses/{id} [put]
func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	var req dto.ExpenseUpdateRequest
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
	if req.Kind != nil {
		expense.Kind = dbmodel.ExpenseKind(*req.Kind)
	}
	if req.Currency != nil {
		expense.Currency = *req.Currency
	}
	if req.Type != nil {
		expense.Type = strings.ToLower(strings.TrimSpace(*req.Type))
	}
	if req.Resource != nil {
		expense.Resource = dbmodel.ExpenseResource(*req.Resource)
	}
	if req.Description != nil {
		expense.Description = *req.Description
	}
	if req.Date != nil {
		// Validate date format
		if _, err := time.Parse("2006-01-02", *req.Date); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, expected YYYY-MM-DD"})
			return
		}
		expense.Date = *req.Date
	}

	if err := h.Service.UpdateExpense(expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toExpenseResponse(expense))
}

// ListExpenses godoc
// @Summary List expenses
// @Description Get list of expenses with optional filtering and ordering
// @Tags expenses
// @Accept json
// @Produce json
// @Param user_id query int false "User ID"
// @Param kind query string false "Expense kind (expense/income)"
// @Param types query []string false "Expense types filter (food/salary/transport/entertainment) - accepts multiple"
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Param order_by query string false "Column to order by: date, amount, type, kind, currency, created_at (default: date)"
// @Param order_dir query string false "Order direction: asc or desc (default: desc)"
// @Success 200 {object} ExpenseListResponse "List of expenses with count"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /expenses [get]
func (h *ExpenseHandler) ListExpenses(c *gin.Context) {
	authCtx := auth.GetAuthContext(c)
	userID := authCtx.UserID
	role := authCtx.Role
	var filter dto.ExpenseFilter
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

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 25
	}

	c.JSON(http.StatusOK, dto.ExpenseListResponse{
		Page:     page,
		PageSize: pageSize,
		Count:    len(expenses),
		Data:     toExpenseResponseList(expenses),
	})
}

// Summary godoc
// @Summary Get expense summary
// @Description Get totals (income, expense, balance) for a filtered set of records.
// @Tags expenses
// @Accept json
// @Produce json
// @Param user_id query int false "User ID"
// @Param kind query string false "Filter by kind (expense/income)"
// @Param types query []string false "Filter by types"
// @Param currencies query []string false "Filter by currencies"
// @Param original_currency query string false "Currency to express totals in (default: VND)"
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} dto.ExpenseSummary "Expense summary"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /expenses/summary [get]
func (h *ExpenseHandler) Summary(c *gin.Context) {
	authCtx := auth.GetAuthContext(c)
	userID := authCtx.UserID
	role := authCtx.Role
	var filter dto.SummaryFilter
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

// Groups godoc
// @Summary Get expense groups
// @Description Get paginated time-bucket groups (DAY/WEEK/MONTH/YEAR) for a filtered set of records.
// @Tags expenses
// @Accept json
// @Produce json
// @Param user_id query int false "User ID"
// @Param kind query string false "Filter by kind (expense/income)"
// @Param types query []string false "Filter by types"
// @Param currencies query []string false "Filter by currencies"
// @Param original_currency query string false "Currency to express totals in (default: VND)"
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Param group_by query string true "Bucket size: DAY, WEEK, MONTH or YEAR"
// @Param page query int false "Page (default: 1)"
// @Param page_size query int false "Page size (default: all)"
// @Success 200 {object} dto.ExpenseGroupsResponse "Paginated expense groups"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /expenses/groups [get]
func (h *ExpenseHandler) Groups(c *gin.Context) {
	authCtx := auth.GetAuthContext(c)
	userID := authCtx.UserID
	role := authCtx.Role
	var filter dto.GroupsFilter
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

	result, err := h.Service.Groups(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch groups"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// DeleteExpense godoc
// @Summary Delete expense
// @Description Delete an expense by ID (user can only delete their own expenses)
// @Tags expenses
// @Accept json
// @Produce json
// @Param id path int true "Expense ID"
// @Success 200 {object} map[string]interface{} "Expense deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid expense ID"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Expense not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /expenses/{id} [delete]
func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	authCtx := auth.GetAuthContext(c)
	userID := authCtx.UserID

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}

	// Check if expense exists and belongs to user
	expense, err := h.Service.GetExpenseByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	role := authCtx.Role
	if role == auth.RoleUser && expense.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own expenses"})
		return
	}

	// Delete expense
	if err := h.Service.DeleteExpense(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expense"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}

// GetUniqueTypes godoc
// @Summary Get unique expense types
// @Description Get list of unique expense type values for the current user
// @Tags expenses
// @Produce json
// @Success 200 {array} string "List of unique types"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /expenses/types [get]
func (h *ExpenseHandler) GetUniqueTypes(c *gin.Context) {
	authCtx := auth.GetAuthContext(c)
	types, err := h.Service.GetUniqueTypes(authCtx.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch types"})
		return
	}
	c.JSON(http.StatusOK, types)
}

// Currency-related endpoints moved to currency package.
