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

func (s *ExpenseService) ListExpenses(userID uint, kind, typeStr string) ([]Expense, error) {
	return s.Repo.ListByUser(userID, kind, typeStr)
}

func (s *ExpenseService) SummaryByDay(userID uint, day time.Time, kind, typeStr string) (float64, error) {
	return s.Repo.SumByDay(userID, day, kind, typeStr)
}
