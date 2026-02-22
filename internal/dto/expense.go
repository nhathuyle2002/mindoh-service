package dto

// ExpenseResponse is the public-facing representation of an expense record.
type ExpenseResponse struct {
	ID          uint    `json:"id"`
	UserID      uint    `json:"user_id"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Kind        string  `json:"kind"`
	Type        string  `json:"type"`
	Resource    string  `json:"resource"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

// ExpenseListResponse wraps the list response with count metadata.
type ExpenseListResponse struct {
	Count int               `json:"count"`
	Data  []ExpenseResponse `json:"data"`
}

// ExpenseCreateRequest is the request body for creating an expense.
type ExpenseCreateRequest struct {
	UserID      uint    `json:"user_id"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Kind        string  `json:"kind"`
	Type        string  `json:"type"`
	Resource    string  `json:"resource"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

// ExpenseUpdateRequest is the request body for updating an expense (all fields optional).
type ExpenseUpdateRequest struct {
	Amount      *float64 `json:"amount,omitempty"`
	Currency    *string  `json:"currency,omitempty"`
	Kind        *string  `json:"kind,omitempty"`
	Type        *string  `json:"type,omitempty"`
	Resource    *string  `json:"resource,omitempty"`
	Description *string  `json:"description,omitempty"`
	Date        *string  `json:"date,omitempty"`
}

// ExpenseFilter holds query parameters for filtering and ordering the expense list.
type ExpenseFilter struct {
	UserID     uint     `form:"user_id"    json:"user_id"`
	Kind       string   `form:"kind"       json:"kind"`
	Type       string   `form:"type"       json:"type"`
	Currencies []string `form:"currencies" json:"currencies"`
	From       string   `form:"from"       json:"from"`
	To         string   `form:"to"         json:"to"`
	OrderBy    string   `form:"order_by"   json:"order_by"`
	OrderDir   string   `form:"order_dir"  json:"order_dir"`
}

// SummaryFilter holds query parameters for the summary endpoint.
type SummaryFilter struct {
	UserID           uint   `form:"user_id"           json:"user_id"`
	OriginalCurrency string `form:"original_currency" json:"original_currency"`
	From             string `form:"from"              json:"from"`
	To               string `form:"to"                json:"to"`
	GroupBy          string `form:"group_by"          json:"group_by"`
}

// CurrencySummary holds per-native-currency income/expense/balance totals.
type CurrencySummary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	TotalBalance float64 `json:"total_balance"`
}

// ExpenseGroup represents aggregated totals for a time bucket (day/month/year).
type ExpenseGroup struct {
	Key         string             `json:"key"`
	Label       string             `json:"label"`
	Income      float64            `json:"income"`
	Expense     float64            `json:"expense"`
	Balance     float64            `json:"balance"`
	TotalByType map[string]float64 `json:"total_by_type"`
}

// ExpenseSummary is the response returned by the summary endpoint.
type ExpenseSummary struct {
	Currency           string                      `json:"currency"`
	TotalIncome        float64                     `json:"total_income"`
	TotalExpense       float64                     `json:"total_expense"`
	TotalBalance       float64                     `json:"total_balance"`
	TotalByTypeIncome  map[string]float64          `json:"total_by_type_income"`
	TotalByTypeExpense map[string]float64          `json:"total_by_type_expense"`
	ByCurrency         map[string]*CurrencySummary `json:"by_currency,omitempty"`
	Groups             []ExpenseGroup              `json:"groups,omitempty"`
}
