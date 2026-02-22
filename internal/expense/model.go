package expense

import dbmodel "mindoh-service/internal/db"

// ExpenseListResponse wraps the list response with count metadata.
type ExpenseListResponse struct {
	Count int               `json:"count"`
	Data  []dbmodel.Expense `json:"data"`
}
