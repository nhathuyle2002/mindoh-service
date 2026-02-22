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

// ExpenseListResponse is a simple paginated list — no aggregated meta.
// Use GET /expenses/summary for totals and GET /expenses/groups for time-bucket data.
type ExpenseListResponse struct {
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Count    int               `json:"count"`
	Data     []ExpenseResponse `json:"data"`
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

// ExpenseFilter holds query parameters for filtering and ordering the paginated expense list.
type ExpenseFilter struct {
	UserID     uint     `form:"user_id"    json:"user_id"`
	Kind       string   `form:"kind"       json:"kind"`
	Types      []string `form:"types"      json:"types"`
	Currencies []string `form:"currencies" json:"currencies"`
	From       string   `form:"from"       json:"from"`
	To         string   `form:"to"         json:"to"`
	OrderBy    string   `form:"order_by"   json:"order_by"`
	OrderDir   string   `form:"order_dir"  json:"order_dir"`
	Page       int      `form:"page"       json:"page"`
	PageSize   int      `form:"page_size"  json:"page_size"`
}

// SummaryFilter holds query parameters for the summary endpoint.
// It supports the same field filters as ExpenseFilter so totals reflect filtered data.
type SummaryFilter struct {
	UserID           uint     `form:"user_id"           json:"user_id"`
	Kind             string   `form:"kind"              json:"kind"`
	Types            []string `form:"types"             json:"types"`
	Currencies       []string `form:"currencies"        json:"currencies"`
	OriginalCurrency string   `form:"original_currency" json:"original_currency"`
	From             string   `form:"from"              json:"from"`
	To               string   `form:"to"                json:"to"`
}

// GroupsFilter holds query parameters for the groups (time-bucket) endpoint.
type GroupsFilter struct {
	UserID           uint     `form:"user_id"           json:"user_id"`
	Kind             string   `form:"kind"              json:"kind"`
	Types            []string `form:"types"             json:"types"`
	Currencies       []string `form:"currencies"        json:"currencies"`
	OriginalCurrency string   `form:"original_currency" json:"original_currency"`
	From             string   `form:"from"              json:"from"`
	To               string   `form:"to"                json:"to"`
	GroupBy          string   `form:"group_by"          json:"group_by"`
	Page             int      `form:"page"              json:"page"`
	PageSize         int      `form:"page_size"         json:"page_size"`
}

// CurrencySummary holds per-native-currency income/expense/balance totals.
type CurrencySummary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	TotalBalance float64 `json:"total_balance"`
}

// ExpenseGroup represents aggregated totals for a time bucket (day/week/month/year).
type ExpenseGroup struct {
	Key         string             `json:"key"`
	Label       string             `json:"label"`
	Income      float64            `json:"income"`
	Expense     float64            `json:"expense"`
	Balance     float64            `json:"balance"`
	TotalByType map[string]float64 `json:"total_by_type"`
}

// ExpenseGroupsResponse is the paginated response from GET /expenses/groups.
type ExpenseGroupsResponse struct {
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
	Groups   []ExpenseGroup `json:"groups"`
}

// ExpenseSummary is the response from GET /expenses/summary.
// Groups have been moved to GET /expenses/groups.
type ExpenseSummary struct {
	Currency           string                      `json:"currency"`
	IncomeCount        int                         `json:"income_count"`
	ExpenseCount       int                         `json:"expense_count"`
	TotalIncome        float64                     `json:"total_income"`
	TotalExpense       float64                     `json:"total_expense"`
	TotalBalance       float64                     `json:"total_balance"`
	TotalByTypeIncome  map[string]float64          `json:"total_by_type_income"`
	TotalByTypeExpense map[string]float64          `json:"total_by_type_expense"`
	ByCurrency         map[string]*CurrencySummary `json:"by_currency,omitempty"`
}
