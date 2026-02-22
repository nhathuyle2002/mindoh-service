package expense

import (
	"strings"

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

func (r *ExpenseRepository) Delete(id uint) error {
	return r.DB.Delete(&Expense{}, id).Error
}

func (r *ExpenseRepository) GetUniqueTypes(userID uint) ([]string, error) {
	var types []string
	db := r.DB.Model(&Expense{}).Distinct("type").Where("type != ''").Order("type asc")
	if userID != 0 {
		db = db.Where("user_id = ?", userID)
	}
	err := db.Pluck("type", &types).Error
	return types, err
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
		db = db.Where("type = ?", strings.ToLower(strings.TrimSpace(filter.Type)))
	}
	if len(filter.Currencies) > 0 {
		db = db.Where("currency IN ?", filter.Currencies)
	}
	if filter.From != "" {
		db = db.Where("date >= ?", filter.From)
	}
	if filter.To != "" {
		db = db.Where("date <= ?", filter.To)
	}
	err := db.Order("date desc").Find(&expenses).Error
	return expenses, err
}

func (r *ExpenseRepository) ListByDateRange(userID uint, from, to string) ([]Expense, error) {
	var expenses []Expense
	db := r.DB.Model(&Expense{})
	if userID != 0 {
		db = db.Where("user_id = ?", userID)
	}
	if from != "" {
		db = db.Where("date >= ?", from)
	}
	if to != "" {
		db = db.Where("date <= ?", to)
	}
	err := db.Order("date desc").Find(&expenses).Error
	return expenses, err
}
