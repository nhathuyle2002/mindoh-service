package expense

import (
	"time"
)

type ExpenseKind string

const (
	ExpenseKindExpense ExpenseKind = "expense"
	ExpenseKindIncome  ExpenseKind = "income"
)

type ExpenseType string

const (
	ExpenseTypeFood          ExpenseType = "food"
	ExpenseTypeSalary        ExpenseType = "salary"
	ExpenseTypeTransport     ExpenseType = "transport"
	ExpenseTypeEntertainment ExpenseType = "entertainment"
)

type Expense struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	UserID      uint        `gorm:"not null" json:"user_id"`
	Amount      float64     `gorm:"not null" json:"amount"`
	Currency    string      `gorm:"type:varchar(3);not null" json:"currency"`
	Kind        ExpenseKind `gorm:"type:varchar(32);not null" json:"kind"` // expense or income
	Type        ExpenseType `gorm:"type:varchar(32);not null" json:"type"` // e.g., food, salary, etc.
	Description string      `gorm:"type:text" json:"description"`
	Date        time.Time   `gorm:"not null" json:"date"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type ExpenseCreateRequest struct {
	UserID      uint        `json:"user_id"`
	Amount      float64     `json:"amount"`
	Currency    string      `json:"currency"` // e.g., USD, EUR
	Kind        ExpenseKind `json:"kind"`     // expense or income
	Type        ExpenseType `json:"type"`     // e.g., food, salary, etc.
	Description string      `json:"description"`
	Date        time.Time   `json:"date"` // Date is optional, if not provided, current time will be used
}
