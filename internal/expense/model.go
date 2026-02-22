package expense

import (
	"time"

	"gorm.io/gorm"
)

type ExpenseKind string

const (
	ExpenseKindExpense ExpenseKind = "expense"
	ExpenseKindIncome  ExpenseKind = "income"
)

type ExpenseResource string

const (
	ExpenseResourceCash   ExpenseResource = "CASH"
	ExpenseResourceCake   ExpenseResource = "CAKE"
	ExpenseResourceVCB    ExpenseResource = "VCB"
	ExpenseResourceVPBank ExpenseResource = "VPBANK"
	ExpenseResourceBIDV   ExpenseResource = "BIDV"
	ExpenseResourceOther  ExpenseResource = "OTHER"
)

type Expense struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	UserID      uint            `gorm:"not null" json:"user_id"`
	Amount      float64         `gorm:"not null" json:"amount"`
	Currency    string          `gorm:"type:varchar(3);not null" json:"currency"`
	Kind        ExpenseKind     `gorm:"type:varchar(32);not null" json:"kind"` // expense or income
	Type        string          `gorm:"type:varchar(32);not null" json:"type"` // e.g., food, salary, etc.
	Resource    ExpenseResource `gorm:"type:varchar(32)" json:"resource"`
	Description string          `gorm:"type:text" json:"description"`
	Date        string          `gorm:"type:varchar(10);not null" json:"date"` // Format: YYYY-MM-DD
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
}

type ExpenseCreateRequest struct {
	UserID      uint            `json:"user_id"`
	Amount      float64         `json:"amount"`
	Currency    string          `json:"currency"`
	Kind        ExpenseKind     `json:"kind"`
	Type        string          `json:"type"`
	Resource    ExpenseResource `json:"resource"`
	Description string          `json:"description"`
	Date        string          `json:"date"` // Format: YYYY-MM-DD
}

type ExpenseFilter struct {
	UserID     uint     `form:"user_id" json:"user_id"`
	Kind       string   `form:"kind" json:"kind"`
	Type       string   `form:"type" json:"type"`
	Currencies []string `form:"currencies" json:"currencies"` // Filter by multiple currencies
	From       string   `form:"from" json:"from"`             // Format: YYYY-MM-DD
	To         string   `form:"to" json:"to"`                 // Format: YYYY-MM-DD
	OrderBy    string   `form:"order_by" json:"order_by"`     // Column: date, amount, type, kind, currency, created_at (default: date)
	OrderDir   string   `form:"order_dir" json:"order_dir"`   // asc or desc (default: desc)
}

// ExpenseListResponse wraps the list response with count metadata
type ExpenseListResponse struct {
	Count int       `json:"count"`
	Data  []Expense `json:"data"`
}

// SummaryFilter is used exclusively for the summary endpoint.
type SummaryFilter struct {
	UserID           uint   `form:"user_id" json:"user_id"`
	OriginalCurrency string `form:"original_currency" json:"original_currency"` // Currency to express totals in (default: VND)
	From             string `form:"from" json:"from"`                           // Format: YYYY-MM-DD
	To               string `form:"to" json:"to"`                               // Format: YYYY-MM-DD
	GroupBy          string `form:"group_by" json:"group_by"`                   // DAY, MONTH, YEAR
}

// CurrencySummary holds per-native-currency income/expense/balance (not converted)
type CurrencySummary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	TotalBalance float64 `json:"total_balance"`
}

type ExpenseSummary struct {
	Currency           string                      `json:"currency"` // Currency totals are expressed in
	TotalIncome        float64                     `json:"total_income"`
	TotalExpense       float64                     `json:"total_expense"`
	TotalBalance       float64                     `json:"total_balance"`
	TotalByTypeIncome  map[string]float64          `json:"total_by_type_income"`  // Income totals per type (converted)
	TotalByTypeExpense map[string]float64          `json:"total_by_type_expense"` // Expense totals per type (absolute, converted)
	ByCurrency         map[string]*CurrencySummary `json:"by_currency,omitempty"` // Per-currency breakdown in native amounts
	Groups             []ExpenseGroup              `json:"groups,omitempty"`      // Only present when group_by is set
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
	Amount      *float64         `json:"amount,omitempty"`
	Currency    *string          `json:"currency,omitempty"`
	Kind        *ExpenseKind     `json:"kind,omitempty"`
	Type        *string          `json:"type,omitempty"`
	Resource    *ExpenseResource `json:"resource,omitempty"`
	Description *string          `json:"description,omitempty"`
	Date        *string          `json:"date,omitempty"` // Format: YYYY-MM-DD
}
