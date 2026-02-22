package expense

import (
	"strings"

	dbmodel "mindoh-service/internal/db"
	"mindoh-service/internal/dto"

	"gorm.io/gorm"
)

// ExpenseRepository handles DB operations for expenses
type ExpenseRepository struct {
	DB *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{DB: db}
}

func (r *ExpenseRepository) GetByID(id uint) (*dbmodel.Expense, error) {
	var expense dbmodel.Expense
	err := r.DB.First(&expense, id).Error
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *ExpenseRepository) Create(expense *dbmodel.Expense) error {
	return r.DB.Create(expense).Error
}

func (r *ExpenseRepository) Update(expense *dbmodel.Expense) error {
	return r.DB.Save(expense).Error
}

func (r *ExpenseRepository) Delete(id uint) error {
	return r.DB.Delete(&dbmodel.Expense{}, id).Error
}

func (r *ExpenseRepository) GetUniqueTypes(userID uint) ([]string, error) {
	var types []string
	query := r.DB.Model(&dbmodel.Expense{}).Distinct("type").Where("type != ''").Order("type asc")
	if userID != 0 {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Pluck("type", &types).Error
	return types, err
}

func (r *ExpenseRepository) ListByFilter(filter dto.ExpenseFilter) ([]dbmodel.Expense, error) {
	var expenses []dbmodel.Expense
	query := r.DB.Model(&dbmodel.Expense{})
	if filter.UserID != 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.Kind != "" {
		query = query.Where("kind = ?", filter.Kind)
	}
	if len(filter.Types) > 0 {
		query = query.Where("type IN ?", filter.Types)
	}
	if len(filter.Currencies) > 0 {
		query = query.Where("currency IN ?", filter.Currencies)
	}
	if filter.From != "" {
		query = query.Where("date >= ?", filter.From)
	}
	if filter.To != "" {
		query = query.Where("date <= ?", filter.To)
	}

	// Build ORDER BY clause
	allowedColumns := map[string]string{
		"date":       "date",
		"amount":     "amount",
		"type":       "type",
		"kind":       "kind",
		"currency":   "currency",
		"created_at": "created_at",
	}
	orderCol := "date"
	if col, ok := allowedColumns[strings.ToLower(filter.OrderBy)]; ok {
		orderCol = col
	}
	orderDir := "desc"
	if strings.ToLower(filter.OrderDir) == "asc" {
		orderDir = "asc"
	}
	query = query.Order(orderCol + " " + orderDir)

	err := query.Find(&expenses).Error
	return expenses, err
}

func (r *ExpenseRepository) ListByDateRange(userID uint, from, to string) ([]dbmodel.Expense, error) {
	var expenses []dbmodel.Expense
	query := r.DB.Model(&dbmodel.Expense{})
	if userID != 0 {
		query = query.Where("user_id = ?", userID)
	}
	if from != "" {
		query = query.Where("date >= ?", from)
	}
	if to != "" {
		query = query.Where("date <= ?", to)
	}
	err := query.Order("date desc").Find(&expenses).Error
	return expenses, err
}
