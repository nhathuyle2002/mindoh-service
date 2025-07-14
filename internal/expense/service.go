package expense

import (
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

func (s *ExpenseService) ListExpenses(filter ExpenseFilter) ([]Expense, error) {
	return s.Repo.ListByFilter(filter)
}

func (s *ExpenseService) SummaryByDay(userID uint, day time.Time, kind, typeStr string) (float64, error) {
	return s.Repo.SumByDay(userID, day, kind, typeStr)
}
