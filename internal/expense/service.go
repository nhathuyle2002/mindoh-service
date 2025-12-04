package expense

import (
	"mindoh-service/internal/currency"
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

func (s *ExpenseService) AddExpense(expense *Expense) error {
	return s.Repo.Create(expense)
}

func (s *ExpenseService) UpdateExpense(expense *Expense) error {
	return s.Repo.Update(expense)
}

func (s *ExpenseService) GetExpenseByID(id uint) (*Expense, error) {
	return s.Repo.GetByID(id)
}

func (s *ExpenseService) DeleteExpense(id uint) error {
	return s.Repo.Delete(id)
}

func (s *ExpenseService) ListExpenses(filter ExpenseFilter) ([]Expense, error) {
	return s.Repo.ListByFilter(filter)
}

func (s *ExpenseService) GetExchangeRates() map[string]float64 {
	// Deprecated: exchange rate access moved to currency package endpoints.
	return nil // This method is now obsolete
}

func (s *ExpenseService) Summary(filter ExpenseFilter) (*ExpenseSummary, error) {
	expenses, err := s.Repo.ListByFilter(filter)
	if err != nil {
		return &ExpenseSummary{}, err
	}

	originalCurrency := filter.OriginalCurrency
	if originalCurrency == "" {
		originalCurrency = "VND"
	}

	var summary *ExpenseSummary
	if len(filter.Currencies) > 0 {
		summary = s.summarizeForCurrencies(expenses, filter.Currencies)
	} else {
		summary = s.summarizeWithConversion(expenses, originalCurrency)
	}

	if filter.GroupBy != "" {
		summary.Groups = s.groupExpenses(expenses, filter, summary.Currency, len(filter.Currencies) == 0)
	}
	return summary, nil
}

func (s *ExpenseService) summarizeForCurrencies(expenses []Expense, currencies []string) *ExpenseSummary {
	byCurrency := make(map[string]*CurrencySummary)

	for _, expense := range expenses {
		// Track by currency
		if byCurrency[expense.Currency] == nil {
			byCurrency[expense.Currency] = &CurrencySummary{
				TotalByType: make(map[string]float64),
			}
		}
		if expense.Kind == ExpenseKindIncome {
			byCurrency[expense.Currency].TotalIncome += expense.Amount
		} else {
			byCurrency[expense.Currency].TotalExpense += expense.Amount
		}
		byCurrency[expense.Currency].Balance = byCurrency[expense.Currency].TotalIncome - byCurrency[expense.Currency].TotalExpense
		byCurrency[expense.Currency].TotalByType[expense.Type] += expense.Amount
	}

	return &ExpenseSummary{
		Expenses:     expenses,
		Currency:     "", // No single currency when filtering by multiple
		TotalIncome:  0,
		TotalExpense: 0,
		Balance:      0,
		TotalByType:  make(map[string]float64),
		ByCurrency:   byCurrency,
	}
}

func (s *ExpenseService) summarizeWithConversion(expenses []Expense, originalCurrency string) *ExpenseSummary {
	var totalIncome, totalExpense float64
	totalByType := make(map[string]float64)
	byCurrency := make(map[string]*CurrencySummary)

	// Get current exchange rates
	exchangeRates := currency.GetExchangeRateService().GetRates()

	// Get exchange rate for original currency
	targetRate := exchangeRates[originalCurrency]
	if targetRate == 0 {
		targetRate = 1
	}

	for _, expense := range expenses {
		// Get exchange rate, default to 1 if not found
		exchangeRate := exchangeRates[expense.Currency]
		if exchangeRate == 0 {
			exchangeRate = 1
		}
		// Convert to original currency (VND rate / target rate)
		convertedAmount := expense.Amount * exchangeRate / targetRate

		// Add to converted totals (in VND)
		if expense.Kind == ExpenseKindIncome {
			totalIncome += convertedAmount
		} else {
			totalExpense += convertedAmount
		}
		totalByType[expense.Type] += convertedAmount

		// Track by original currency
		if byCurrency[expense.Currency] == nil {
			byCurrency[expense.Currency] = &CurrencySummary{
				TotalByType: make(map[string]float64),
			}
		}
		if expense.Kind == ExpenseKindIncome {
			byCurrency[expense.Currency].TotalIncome += expense.Amount
		} else {
			byCurrency[expense.Currency].TotalExpense += expense.Amount
		}
		byCurrency[expense.Currency].Balance = byCurrency[expense.Currency].TotalIncome - byCurrency[expense.Currency].TotalExpense
		byCurrency[expense.Currency].TotalByType[expense.Type] += expense.Amount
	}

	return &ExpenseSummary{
		Expenses:     expenses,
		Currency:     originalCurrency, // Converted to selected original currency
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		Balance:      totalIncome - totalExpense,
		TotalByType:  totalByType,
		ByCurrency:   byCurrency,
	}
}

// groupExpenses aggregates expenses according to filter.GroupBy (DAY, MONTH, YEAR).
// When convertTotals is true (no explicit currencies filter) amounts are converted into targetCurrency.
func (s *ExpenseService) groupExpenses(expenses []Expense, filter ExpenseFilter, targetCurrency string, convertTotals bool) []ExpenseGroup {
	if len(expenses) == 0 {
		return []ExpenseGroup{}
	}
	mode := strings.ToUpper(filter.GroupBy)
	if mode != "DAY" && mode != "MONTH" && mode != "YEAR" {
		return []ExpenseGroup{}
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
			key = exp.Date.Format("2006-01-02")
		case "MONTH":
			key = exp.Date.Format("2006-01")
		case "YEAR":
			key = exp.Date.Format("2006")
		}
		if groups[key] == nil {
			groups[key] = &agg{TotalByType: make(map[string]float64)}
		}
		amount := exp.Amount
		if convertTotals {
			rate := exchangeRates[exp.Currency]
			if rate == 0 {
				rate = 1
			}
			amount = amount * rate / targetRate
		}
		if exp.Kind == ExpenseKindIncome {
			groups[key].Income += amount
		} else {
			groups[key].Expense += amount
		}
		groups[key].TotalByType[exp.Type] += amount
	}
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		// Parse key back to time for ordering
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
	result := make([]ExpenseGroup, 0, len(keys))
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
		case "YEAR":
			label = k
		}
		result = append(result, ExpenseGroup{
			Key:         k,
			Label:       label,
			Income:      g.Income,
			Expense:     g.Expense,
			Balance:     g.Income - g.Expense,
			TotalByType: g.TotalByType,
		})
	}
	return result
}
