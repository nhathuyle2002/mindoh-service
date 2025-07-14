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

func (s *ExpenseService) Summary(filter ExpenseFilter) (*ExpenseSummary, error) {
	expenses, err := s.Repo.ListByFilter(filter)
	if err != nil {
		return &ExpenseSummary{}, err
	}
	totalByKind := make(map[ExpenseKind]float64)
	totalByType := make(map[ExpenseType]float64)
	totalAmount := 0.0
	for _, expense := range expenses {
		totalByKind[expense.Kind] += expense.Amount
		totalByType[expense.Type] += expense.Amount
		totalAmount += expense.Amount
	}
	return &ExpenseSummary{
		Expenses:    expenses,
		TotalByKind: totalByKind,
		TotalByType: totalByType,
		TotalAmount: totalAmount,
	}, nil
}
