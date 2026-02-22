package db

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

// Expense is the database model for an expense or income record.
type Expense struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	UserID      uint            `gorm:"not null" json:"user_id"`
	Amount      float64         `gorm:"not null" json:"amount"`
	Currency    string          `gorm:"type:varchar(3);not null" json:"currency"`
	Kind        ExpenseKind     `gorm:"type:varchar(32);not null" json:"kind"`
	Type        string          `gorm:"type:varchar(32);not null" json:"type"`
	Resource    ExpenseResource `gorm:"type:varchar(32)" json:"resource"`
	Description string          `gorm:"type:text" json:"description"`
	Date        string          `gorm:"type:varchar(10);not null" json:"date"` // Format: YYYY-MM-DD
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
}
