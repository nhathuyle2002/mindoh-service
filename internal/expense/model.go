package expense

import (
	"time"
)

type ExpenseKind string

const (
	ExpenseKindExpense ExpenseKind = "expense"
	ExpenseKindIncome  ExpenseKind = "income"
)

type Expense struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	UserID      uint        `gorm:"not null" json:"user_id"`
	Amount      float64     `gorm:"not null" json:"amount"`
	Currency    string      `gorm:"type:varchar(3);not null" json:"currency"`
	Kind        ExpenseKind `gorm:"type:varchar(32);not null" json:"kind"` // expense or income
	Type        string      `gorm:"type:varchar(32);not null" json:"type"` // e.g., food, salary, etc.
	Description string      `gorm:"type:text" json:"description"`
	Date        time.Time   `gorm:"not null" json:"date"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type ExpenseCreateRequest struct {
	UserID      uint        `json:"user_id"`
	Amount      float64     `json:"amount"`
	Currency    string      `json:"currency"`
	Kind        ExpenseKind `json:"kind"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Date        time.Time   `json:"date"` // Date is optional, if not provided, current time will be used
}

type ExpenseFilter struct {
	UserID           uint      `form:"user_id" json:"user_id"`
	Kind             string    `form:"kind" json:"kind"`
	Type             string    `form:"type" json:"type"`
	Currencies       []string  `form:"currencies" json:"currencies"`               // Filter by multiple currencies
	OriginalCurrency string    `form:"original_currency" json:"original_currency"` // Currency to convert totals into when no currency filter
	From             time.Time `form:"from" json:"from"`
	To               time.Time `form:"to" json:"to"`
	GroupBy          string    `form:"group_by" json:"group_by"` // DAY, MONTH, YEAR (uppercase, case-insensitive)
}

type ExpenseSummary struct {
	Expenses     []Expense                   `json:"expenses"`
	Currency     string                      `json:"currency"` // The currency for totals (if filtered) or base currency (VND)
	TotalIncome  float64                     `json:"total_income"`
	TotalExpense float64                     `json:"total_expense"`
	Balance      float64                     `json:"balance"`
	TotalByType  map[string]float64          `json:"total_by_type"`
	ByCurrency   map[string]*CurrencySummary `json:"by_currency,omitempty"` // Only present when no currency filter
	Groups       []ExpenseGroup              `json:"groups,omitempty"`      // Optional grouped summary buckets
}

type CurrencySummary struct {
	TotalIncome  float64            `json:"total_income"`
	TotalExpense float64            `json:"total_expense"`
	Balance      float64            `json:"balance"`
	TotalByType  map[string]float64 `json:"total_by_type"`
}

// ExpenseGroup represents aggregated totals for a grouping bucket (day/month/year)
type ExpenseGroup struct {
	Key         string             `json:"key"`   // Raw key (YYYY-MM-DD / YYYY-MM / YYYY)
	Label       string             `json:"label"` // Human-friendly label
	Income      float64            `json:"income"`
	Expense     float64            `json:"expense"`
	Balance     float64            `json:"balance"`
	TotalByType map[string]float64 `json:"total_by_type"`
}

type ExpenseUpdateRequest struct {
	Amount      *float64     `json:"amount,omitempty"`
	Currency    *string      `json:"currency,omitempty"`
	Kind        *ExpenseKind `json:"kind,omitempty"`
	Type        *string      `json:"type,omitempty"`
	Description *string      `json:"description,omitempty"`
	Date        *time.Time   `json:"date,omitempty"`
}
