package expense

import (
	"fmt"
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

func (r *ExpenseRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.DB.Model(&dbmodel.Expense{}).Where("id = ?", id).Updates(fields).Error
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

// buildBaseQuery applies all WHERE clauses from filter (no ORDER BY, no LIMIT).
func (r *ExpenseRepository) buildBaseQuery(filter dto.ExpenseFilter) *gorm.DB {
	q := r.DB.Model(&dbmodel.Expense{})
	if filter.UserID != 0 {
		q = q.Where("user_id = ?", filter.UserID)
	}
	if filter.Kind != "" {
		q = q.Where("kind = ?", filter.Kind)
	}
	if len(filter.Types) > 0 {
		q = q.Where("type IN ?", filter.Types)
	}
	if len(filter.Currencies) > 0 {
		q = q.Where("currency IN ?", filter.Currencies)
	}
	if filter.From != "" {
		q = q.Where("date >= ?", filter.From)
	}
	if filter.To != "" {
		q = q.Where("date <= ?", filter.To)
	}
	return q
}

// AggregateMeta runs a single lightweight GROUP BY query to compute
// total count, income/expense counts, and per-currency totals — no row fetching.
func (r *ExpenseRepository) AggregateMeta(filter dto.ExpenseFilter) (total, incomeCount, expenseCount int, byCurrency map[string]*dto.CurrencySummary, err error) {
	type row struct {
		Kind     string  `gorm:"column:kind"`
		Currency string  `gorm:"column:currency"`
		Cnt      int     `gorm:"column:cnt"`
		SumAmt   float64 `gorm:"column:sum_amt"`
	}
	var rows []row
	err = r.buildBaseQuery(filter).
		Select("kind, currency, COUNT(*) AS cnt, SUM(amount) AS sum_amt").
		Group("kind, currency").
		Scan(&rows).Error
	if err != nil {
		return
	}
	byCurrency = map[string]*dto.CurrencySummary{}
	for _, rw := range rows {
		total += rw.Cnt
		if rw.Kind == "income" {
			incomeCount += rw.Cnt
		} else {
			expenseCount += rw.Cnt
		}
		cs, ok := byCurrency[rw.Currency]
		if !ok {
			cs = &dto.CurrencySummary{}
			byCurrency[rw.Currency] = cs
		}
		if rw.Kind == "income" {
			cs.TotalIncome += rw.SumAmt
		} else {
			cs.TotalExpense += rw.SumAmt
		}
		cs.TotalBalance += rw.SumAmt
	}
	return
}

func (r *ExpenseRepository) ListByFilter(filter dto.ExpenseFilter) ([]dbmodel.Expense, error) {
	var expenses []dbmodel.Expense

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

	query := r.buildBaseQuery(filter).Order(orderCol + " " + orderDir)

	// DB-level pagination
	if filter.PageSize > 0 {
		page := filter.Page
		if page < 1 {
			page = 1
		}
		offset := (page - 1) * filter.PageSize
		query = query.Limit(filter.PageSize).Offset(offset)
	}

	err := query.Find(&expenses).Error
	return expenses, err
}

// ListAllByFilter fetches all matching rows with no LIMIT/OFFSET — used by the summary endpoint
// which needs every matching record to compute aggregated totals.
func (r *ExpenseRepository) ListAllByFilter(filter dto.ExpenseFilter) ([]dbmodel.Expense, error) {
	var expenses []dbmodel.Expense
	err := r.buildBaseQuery(filter).Order("date desc").Find(&expenses).Error
	return expenses, err
}

// GroupAggRow is one row returned by ListGroupsAggByFilter.
type GroupAggRow struct {
	Bucket   string  `gorm:"column:bucket"`
	Currency string  `gorm:"column:currency"`
	Type     string  `gorm:"column:type"`
	Kind     string  `gorm:"column:kind"`
	Total    float64 `gorm:"column:total"`
}

// bucketSQL returns the PostgreSQL expression that maps a varchar date (YYYY-MM-DD)
// to a time-bucket key string that sorts lexicographically.
func bucketSQL(groupBy string) (string, error) {
	switch strings.ToUpper(groupBy) {
	case "DAY":
		return "date", nil
	case "WEEK":
		// Monday of the ISO week, kept as YYYY-MM-DD so it sorts naturally.
		return "TO_CHAR(DATE_TRUNC('week', date::date), 'YYYY-MM-DD')", nil
	case "MONTH":
		return "SUBSTRING(date, 1, 7)", nil
	case "YEAR":
		return "SUBSTRING(date, 1, 4)", nil
	default:
		return "", fmt.Errorf("unsupported group_by value: %s", groupBy)
	}
}

// ListGroupsAggByFilter returns all time-bucket aggregation rows for the filter.
// Sorting and pagination are handled in the service layer (Go-side) so that
// computed fields like income/expense/balance (which require live exchange-rate
// conversion) can be sorted accurately.
func (r *ExpenseRepository) ListGroupsAggByFilter(filter dto.GroupsFilter) (rows []GroupAggRow, err error) {
	expr, err := bucketSQL(filter.GroupBy)
	if err != nil {
		return
	}

	lf := dto.ExpenseFilter{
		UserID:     filter.UserID,
		Kind:       filter.Kind,
		Types:      filter.Types,
		Currencies: filter.Currencies,
		From:       filter.From,
		To:         filter.To,
	}

	err = r.buildBaseQuery(lf).
		Select(fmt.Sprintf("%s AS bucket, currency, type, kind, SUM(amount) AS total", expr)).
		Group(fmt.Sprintf("%s, currency, type, kind", expr)).
		Order("bucket DESC").
		Scan(&rows).Error
	return
}
