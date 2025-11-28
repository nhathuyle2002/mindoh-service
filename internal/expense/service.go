package expense

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

func (s *ExpenseService) ListExpenses(filter ExpenseFilter) ([]Expense, error) {
	return s.Repo.ListByFilter(filter)
}

func (s *ExpenseService) GetExchangeRates() map[string]float64 {
	return GetExchangeRateService().GetRates()
}

func (s *ExpenseService) Summary(filter ExpenseFilter) (*ExpenseSummary, error) {
	expenses, err := s.Repo.ListByFilter(filter)
	if err != nil {
		return &ExpenseSummary{}, err
	}

	// Set default currency to VND if not specified
	defaultCurrency := filter.DefaultCurrency
	if defaultCurrency == "" {
		defaultCurrency = "VND"
	}

	// If currencies filter is specified, return summary for those currencies only (no conversion)
	if len(filter.Currencies) > 0 {
		return s.summarizeForCurrencies(expenses, filter.Currencies), nil
	}

	// Otherwise, return converted summary in default currency + breakdown by currency
	return s.summarizeWithConversion(expenses, defaultCurrency), nil
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

func (s *ExpenseService) summarizeWithConversion(expenses []Expense, defaultCurrency string) *ExpenseSummary {
	var totalIncome, totalExpense float64
	totalByType := make(map[string]float64)
	byCurrency := make(map[string]*CurrencySummary)

	// Get current exchange rates
	exchangeRates := GetExchangeRateService().GetRates()

	// Get exchange rate for default currency
	defaultRate := exchangeRates[defaultCurrency]
	if defaultRate == 0 {
		defaultRate = 1
	}

	for _, expense := range expenses {
		// Get exchange rate, default to 1 if not found
		exchangeRate := exchangeRates[expense.Currency]
		if exchangeRate == 0 {
			exchangeRate = 1
		}
		// Convert to default currency (VND rate / default rate)
		convertedAmount := expense.Amount * exchangeRate / defaultRate

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
		Currency:     defaultCurrency, // Converted to default currency
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		Balance:      totalIncome - totalExpense,
		TotalByType:  totalByType,
		ByCurrency:   byCurrency,
	}
}
