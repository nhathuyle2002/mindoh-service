package expense

import (
	"errors"
	"mindoh-service/internal/currency"
	dbmodel "mindoh-service/internal/db"
	"mindoh-service/internal/dto"
	"sort"
	"strings"
	"time"
)

// ExpenseService handles business logic for expenses
type ExpenseService struct {
	Repo *ExpenseRepository
}

func NewExpenseService(repo *ExpenseRepository) *ExpenseService {
	return &ExpenseService{Repo: repo}
}

func (s *ExpenseService) AddExpense(expense *dbmodel.Expense) error {
	// Validate amount sign based on kind
	if expense.Kind == dbmodel.ExpenseKindExpense && expense.Amount > 0 {
		return errors.New("expense amount must be negative")
	}
	if expense.Kind == dbmodel.ExpenseKindIncome && expense.Amount < 0 {
		return errors.New("income amount must be positive")
	}
	return s.Repo.Create(expense)
}

func (s *ExpenseService) UpdateExpense(expense *dbmodel.Expense) error {
	// Validate amount sign based on kind
	if expense.Kind == dbmodel.ExpenseKindExpense && expense.Amount > 0 {
		return errors.New("expense amount must be negative")
	}
	if expense.Kind == dbmodel.ExpenseKindIncome && expense.Amount < 0 {
		return errors.New("income amount must be positive")
	}
	return s.Repo.Update(expense)
}

// UpdateExpenseFields updates only the explicitly provided fields for an expense.
// expense is the current DB state (used for validation of the final kind/amount).
func (s *ExpenseService) UpdateExpenseFields(expense *dbmodel.Expense, fields map[string]interface{}) error {
	// Validate final kind/amount sign
	if expense.Kind == dbmodel.ExpenseKindExpense && expense.Amount > 0 {
		return errors.New("expense amount must be negative")
	}
	if expense.Kind == dbmodel.ExpenseKindIncome && expense.Amount < 0 {
		return errors.New("income amount must be positive")
	}
	return s.Repo.UpdateFields(expense.ID, fields)
}

func (s *ExpenseService) GetExpenseByID(id uint) (*dbmodel.Expense, error) {
	return s.Repo.GetByID(id)
}

func (s *ExpenseService) GetUniqueTypes(userID uint) ([]string, error) {
	return s.Repo.GetUniqueTypes(userID)
}

func (s *ExpenseService) DeleteExpense(id uint) error {
	return s.Repo.Delete(id)
}

func (s *ExpenseService) ListExpenses(filter dto.ExpenseFilter) ([]dbmodel.Expense, error) {
	return s.Repo.ListByFilter(filter)
}

func (s *ExpenseService) AggregateMeta(filter dto.ExpenseFilter) (total, incomeCount, expenseCount int, byCurrency map[string]*dto.CurrencySummary, err error) {
	return s.Repo.AggregateMeta(filter)
}

func (s *ExpenseService) Summary(filter dto.SummaryFilter) (*dto.ExpenseSummary, error) {
	listFilter := dto.ExpenseFilter{
		UserID:     filter.UserID,
		Kind:       filter.Kind,
		Types:      filter.Types,
		Currencies: filter.Currencies,
		From:       filter.From,
		To:         filter.To,
	}
	expenses, err := s.Repo.ListAllByFilter(listFilter)
	if err != nil {
		return &dto.ExpenseSummary{}, err
	}

	originalCurrency := filter.OriginalCurrency
	if originalCurrency == "" {
		originalCurrency = "VND"
	}

	summary := s.computeSummary(expenses, originalCurrency)
	return summary, nil
}

func (s *ExpenseService) computeSummary(expenses []dbmodel.Expense, targetCurrency string) *dto.ExpenseSummary {
	var totalIncome, totalExpense float64
	var incomeCount, expenseCount int
	totalByTypeIncome := make(map[string]float64)
	totalByTypeExpense := make(map[string]float64)
	byCurrency := make(map[string]*dto.CurrencySummary)

	exchangeRates := currency.GetExchangeRateService().GetRates()
	targetRate := exchangeRates[targetCurrency]
	if targetRate == 0 {
		targetRate = 1
	}

	for _, expense := range expenses {
		rate := exchangeRates[expense.Currency]
		if rate == 0 {
			rate = 1
		}
		converted := expense.Amount * rate / targetRate

		// Per-currency native amounts (not converted)
		if _, ok := byCurrency[expense.Currency]; !ok {
			byCurrency[expense.Currency] = &dto.CurrencySummary{}
		}

		if expense.Kind == dbmodel.ExpenseKindIncome {
			incomeCount++
			totalIncome += converted
			totalByTypeIncome[expense.Type] += converted
			byCurrency[expense.Currency].TotalIncome += expense.Amount
		} else {
			expenseCount++
			totalExpense += converted
			totalByTypeExpense[expense.Type] += converted
			byCurrency[expense.Currency].TotalExpense += expense.Amount
		}
	}

	// Compute per-currency balance
	for _, cs := range byCurrency {
		cs.TotalBalance = cs.TotalIncome + cs.TotalExpense
	}

	// Only include ByCurrency when there are multiple currencies
	var byCurrencyResult map[string]*dto.CurrencySummary
	if len(byCurrency) > 1 {
		byCurrencyResult = byCurrency
	}

	return &dto.ExpenseSummary{
		Currency:           targetCurrency,
		IncomeCount:        incomeCount,
		ExpenseCount:       expenseCount,
		TotalIncome:        totalIncome,
		TotalExpense:       totalExpense,
		TotalBalance:       totalIncome + totalExpense,
		TotalByTypeIncome:  totalByTypeIncome,
		TotalByTypeExpense: totalByTypeExpense,
		ByCurrency:         byCurrencyResult,
	}
}

// Groups aggregates expenses into time-bucket groups with Go-side sort and pagination.
// Currency conversion uses live rates; sort/paginate in Go so computed fields are accurate.
func (s *ExpenseService) Groups(filter dto.GroupsFilter) (*dto.ExpenseGroupsResponse, error) {
	aggRows, err := s.Repo.ListGroupsAggByFilter(filter)
	if err != nil {
		return nil, err
	}

	originalCurrency := filter.OriginalCurrency
	if originalCurrency == "" {
		originalCurrency = "VND"
	}

	exchangeRates := currency.GetExchangeRateService().GetRates()
	targetRate := exchangeRates[originalCurrency]
	if targetRate == 0 {
		targetRate = 1
	}

	type agg struct {
		Income      float64
		Expense     float64
		TotalByType map[string]float64
	}
	groupMap := make(map[string]*agg)
	keyOrder := make([]string, 0)

	for _, row := range aggRows {
		if groupMap[row.Bucket] == nil {
			groupMap[row.Bucket] = &agg{TotalByType: make(map[string]float64)}
			keyOrder = append(keyOrder, row.Bucket)
		}
		rate := exchangeRates[row.Currency]
		if rate == 0 {
			rate = 1
		}
		converted := row.Total * rate / targetRate
		groupMap[row.Bucket].TotalByType[row.Type] += converted
		if row.Kind == string(dbmodel.ExpenseKindIncome) {
			groupMap[row.Bucket].Income += converted
		} else {
			groupMap[row.Bucket].Expense += converted
		}
	}

	mode := strings.ToUpper(filter.GroupBy)
	groups := make([]dto.ExpenseGroup, 0, len(keyOrder))
	for _, k := range keyOrder {
		g := groupMap[k]
		groups = append(groups, dto.ExpenseGroup{
			Key:         k,
			Label:       bucketLabel(k, mode),
			Income:      g.Income,
			Expense:     g.Expense,
			Balance:     g.Income + g.Expense,
			TotalByType: g.TotalByType,
		})
	}

	// Sort only when the caller explicitly requests it; otherwise keep
	// the natural DB order (bucket DESC from the repository query).
	orderBy := strings.ToLower(filter.OrderBy)
	if orderBy != "" {
		orderDir := strings.ToLower(filter.OrderDir)
		if orderDir != "asc" {
			orderDir = "desc"
		}
		sort.Slice(groups, func(i, j int) bool {
			var less bool
			switch orderBy {
			case "income":
				less = groups[i].Income < groups[j].Income
			case "expense":
				less = groups[i].Expense < groups[j].Expense
			case "balance":
				less = groups[i].Balance < groups[j].Balance
			default: // "period"
				less = groups[i].Key < groups[j].Key
			}
			if orderDir == "desc" {
				return !less
			}
			return less
		})
	}

	total := len(groups)

	// Paginate.
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	var pagedGroups []dto.ExpenseGroup
	if pageSize > 0 {
		start := (page - 1) * pageSize
		if start >= total {
			pagedGroups = []dto.ExpenseGroup{}
		} else {
			end := start + pageSize
			if end > total {
				end = total
			}
			pagedGroups = groups[start:end]
		}
	} else {
		pagedGroups = groups
	}

	return &dto.ExpenseGroupsResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Groups:   pagedGroups,
	}, nil
}

// bucketLabel converts a bucket key to a human-friendly label.
func bucketLabel(key, mode string) string {
	switch mode {
	case "DAY":
		if t, err := time.Parse("2006-01-02", key); err == nil {
			return t.Format("02 Jan 2006")
		}
	case "WEEK":
		// key is YYYY-MM-DD (Monday of the week)
		if t, err := time.Parse("2006-01-02", key); err == nil {
			return "W/o " + t.Format("02 Jan")
		}
	case "MONTH":
		if t, err := time.Parse("2006-01", key); err == nil {
			return t.Format("Jan 2006")
		}
	case "YEAR":
		return key
	}
	return key
}
