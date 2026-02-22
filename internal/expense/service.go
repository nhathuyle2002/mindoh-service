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

func (s *ExpenseService) Summary(filter dto.SummaryFilter) (*dto.ExpenseSummary, error) {
	expenses, err := s.Repo.ListByDateRange(filter.UserID, filter.From, filter.To)
	if err != nil {
		return &dto.ExpenseSummary{}, err
	}

	originalCurrency := filter.OriginalCurrency
	if originalCurrency == "" {
		originalCurrency = "VND"
	}

	summary := s.computeSummary(expenses, originalCurrency)

	if filter.GroupBy != "" {
		summary.Groups = s.groupExpenses(expenses, filter.GroupBy, originalCurrency)
	}
	return summary, nil
}

func (s *ExpenseService) computeSummary(expenses []dbmodel.Expense, targetCurrency string) *dto.ExpenseSummary {
	var totalIncome, totalExpense float64
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
			totalIncome += converted
			totalByTypeIncome[expense.Type] += converted
			byCurrency[expense.Currency].TotalIncome += expense.Amount
		} else {
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
		TotalIncome:        totalIncome,
		TotalExpense:       totalExpense,
		TotalBalance:       totalIncome + totalExpense,
		TotalByTypeIncome:  totalByTypeIncome,
		TotalByTypeExpense: totalByTypeExpense,
		ByCurrency:         byCurrencyResult,
	}
}

// groupExpenses aggregates expenses into buckets according to groupBy (DAY, MONTH, YEAR),
// converting all amounts to targetCurrency.
func (s *ExpenseService) groupExpenses(expenses []dbmodel.Expense, groupBy string, targetCurrency string) []dto.ExpenseGroup {
	if len(expenses) == 0 {
		return []dto.ExpenseGroup{}
	}
	mode := strings.ToUpper(groupBy)
	if mode != "DAY" && mode != "MONTH" && mode != "YEAR" {
		return []dto.ExpenseGroup{}
	}

	exchangeRates := currency.GetExchangeRateService().GetRates()
	targetRate := exchangeRates[targetCurrency]
	if targetRate == 0 {
		targetRate = 1
	}

	type agg struct {
		Income      float64
		Expense     float64
		TotalByType map[string]float64
	}
	groups := make(map[string]*agg)

	for _, exp := range expenses {
		var key string
		switch mode {
		case "DAY":
			key = exp.Date
		case "MONTH":
			if len(exp.Date) >= 7 {
				key = exp.Date[:7]
			}
		case "YEAR":
			if len(exp.Date) >= 4 {
				key = exp.Date[:4]
			}
		}
		if groups[key] == nil {
			groups[key] = &agg{TotalByType: make(map[string]float64)}
		}
		rate := exchangeRates[exp.Currency]
		if rate == 0 {
			rate = 1
		}
		amount := exp.Amount * rate / targetRate
		groups[key].TotalByType[exp.Type] += amount
		if exp.Kind == dbmodel.ExpenseKindIncome {
			groups[key].Income += amount
		} else {
			groups[key].Expense += amount
		}
	}

	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		var ti, tj time.Time
		switch mode {
		case "DAY":
			ti, _ = time.Parse("2006-01-02", keys[i])
			tj, _ = time.Parse("2006-01-02", keys[j])
		case "MONTH":
			ti, _ = time.Parse("2006-01", keys[i])
			tj, _ = time.Parse("2006-01", keys[j])
		case "YEAR":
			ti, _ = time.Parse("2006", keys[i])
			tj, _ = time.Parse("2006", keys[j])
		}
		return ti.Before(tj)
	})

	result := make([]dto.ExpenseGroup, 0, len(keys))
	for _, k := range keys {
		g := groups[k]
		label := k
		switch mode {
		case "DAY":
			t, _ := time.Parse("2006-01-02", k)
			label = t.Format("02 Jan 2006")
		case "MONTH":
			t, _ := time.Parse("2006-01", k)
			label = t.Format("Jan 2006")
		}
		result = append(result, dto.ExpenseGroup{
			Key:         k,
			Label:       label,
			Income:      g.Income,
			Expense:     g.Expense,
			Balance:     g.Income + g.Expense,
			TotalByType: g.TotalByType,
		})
	}
	return result
}
