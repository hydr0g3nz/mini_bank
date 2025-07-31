// internal/application/dto/mapper.go
package dto

import (
	"github.com/hydr0g3nz/mini_bank/internal/domain/entity"
	"github.com/hydr0g3nz/mini_bank/internal/domain/vo"
)

// AccountMapper provides mapping between Account entity and DTOs
type AccountMapper struct{}

// ToResponse converts Account entity to AccountResponse DTO
func (m *AccountMapper) ToResponse(account *entity.Account) AccountResponse {
	return AccountResponse{
		ID:          account.ID.String(),
		AccountName: account.AccountName,
		Balance:     account.Balance.Amount().InexactFloat64(),
		Status:      string(account.Status),
		CreatedAt:   account.CreatedAt,
		UpdatedAt:   account.UpdatedAt,
	}
}

// ToResponseList converts slice of Account entities to AccountListResponse DTO
func (m *AccountMapper) ToResponseList(accounts []*entity.Account, pagination PaginationInfo) AccountListResponse {
	responses := make([]AccountResponse, len(accounts))
	for i, account := range accounts {
		responses[i] = m.ToResponse(account)
	}

	return AccountListResponse{
		Accounts:   responses,
		Pagination: pagination,
	}
}

// FromCreateRequest converts CreateAccountRequest DTO to domain values
func (m *AccountMapper) FromCreateRequest(req CreateAccountRequest) (string, vo.Money, error) {
	money := vo.NewMoneyFromFloat(req.InitialBalance)
	return req.AccountName, money, nil
}

// TransactionMapper provides mapping between Transaction entity and DTOs
type TransactionMapper struct{}

// ToResponse converts Transaction entity to TransactionResponse DTO
func (m *TransactionMapper) ToResponse(transaction *entity.Transaction) TransactionResponse {
	response := TransactionResponse{
		ID:              transaction.ID.String(),
		TransactionType: string(transaction.TransactionType),
		Amount:          transaction.Amount.Amount().InexactFloat64(),
		Description:     transaction.Description,
		Reference:       transaction.Reference,
		Status:          string(transaction.Status),
		CreatedAt:       transaction.CreatedAt,
		CompletedAt:     transaction.CompletedAt,
	}

	if transaction.FromAccountID != nil {
		fromID := transaction.FromAccountID.String()
		response.FromAccountID = &fromID
	}

	if transaction.ToAccountID != nil {
		toID := transaction.ToAccountID.String()
		response.ToAccountID = &toID
	}

	return response
}

// ToResponseList converts slice of Transaction entities to TransactionListResponse DTO
func (m *TransactionMapper) ToResponseList(transactions []*entity.Transaction, pagination PaginationInfo) TransactionListResponse {
	responses := make([]TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = m.ToResponse(transaction)
	}

	return TransactionListResponse{
		Transactions: responses,
		Pagination:   pagination,
	}
}

// FromCreateRequest converts CreateTransactionRequest DTO to domain values
func (m *TransactionMapper) FromCreateRequest(req CreateTransactionRequest) (
	fromAccountID *vo.AccountID,
	toAccountID *vo.AccountID,
	transactionType vo.TransactionType,
	amount vo.Money,
	description string,
	reference string,
	err error,
) {
	// Parse amount
	amount = vo.NewMoneyFromFloat(req.Amount)

	// Parse transaction type
	transactionType = vo.TransactionType(req.TransactionType)

	// Parse account IDs
	if req.FromAccountID != nil {
		fromID, parseErr := vo.NewAccountIDFromString(*req.FromAccountID)
		if parseErr != nil {
			return nil, nil, "", vo.Money{}, "", "", parseErr
		}
		fromAccountID = &fromID
	}

	if req.ToAccountID != nil {
		toID, parseErr := vo.NewAccountIDFromString(*req.ToAccountID)
		if parseErr != nil {
			return nil, nil, "", vo.Money{}, "", "", parseErr
		}
		toAccountID = &toID
	}

	description = req.Description
	reference = req.Reference

	return fromAccountID, toAccountID, transactionType, amount, description, reference, nil
}
