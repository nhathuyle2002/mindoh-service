package expense

import (
	"gorm.io/gorm"
)

// ExpenseRepository handles DB operations for expenses
type ExpenseRepository struct {
	DB *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{DB: db}
}

func (r *ExpenseRepository) GetByID(id uint) (*Expense, error) {
	var expense Expense
	err := r.DB.First(&expense, id).Error
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *ExpenseRepository) Create(expense *Expense) error {
	return r.DB.Create(expense).Error
}

func (r *ExpenseRepository) Update(expense *Expense) error {
	return r.DB.Save(expense).Error
}

func (r *ExpenseRepository) ListByFilter(filter ExpenseFilter) ([]Expense, error) {
	var expenses []Expense
	db := r.DB.Model(&Expense{})
	if filter.UserID != 0 {
		db = db.Where("user_id = ?", filter.UserID)
	}
	if filter.Kind != "" {
		db = db.Where("kind = ?", filter.Kind)
	}
	if filter.Type != "" {
		db = db.Where("type = ?", filter.Type)
	}
	if filter.Currency != "" {
		db = db.Where("currency = ?", filter.Currency)
	}
	if !filter.From.IsZero() {
		db = db.Where("date >= ?", filter.From)
	}
	if !filter.To.IsZero() {
		db = db.Where("date <= ?", filter.To)
	}
	err := db.Order("date desc").Find(&expenses).Error
	return expenses, err
}
