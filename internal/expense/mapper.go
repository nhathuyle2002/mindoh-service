package expense

import (
	dbmodel "mindoh-service/internal/db"
	"mindoh-service/internal/dto"
)

func toExpenseResponse(e *dbmodel.Expense) dto.ExpenseResponse {
	return dto.ExpenseResponse{
		ID:          e.ID,
		UserID:      e.UserID,
		Amount:      e.Amount,
		Currency:    e.Currency,
		Kind:        string(e.Kind),
		Type:        e.Type,
		Resource:    string(e.Resource),
		Description: e.Description,
		Date:        e.Date,
	}
}

func toExpenseResponseList(expenses []dbmodel.Expense) []dto.ExpenseResponse {
	result := make([]dto.ExpenseResponse, len(expenses))
	for i := range expenses {
		result[i] = toExpenseResponse(&expenses[i])
	}
	return result
}
