package entity

import (
	"strings"
	"time"

	errs "github.com/hydr0g3nz/mini_bank/internal/domain/error"
	"github.com/hydr0g3nz/mini_bank/internal/domain/vo"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID              vo.TransactionID     `json:"id"`
	FromAccountID   *vo.AccountID        `json:"from_account_id,omitempty"`
	ToAccountID     *vo.AccountID        `json:"to_account_id,omitempty"`
	TransactionType vo.TransactionType   `json:"transaction_type"`
	Amount          vo.Money             `json:"amount"`
	Description     string               `json:"description"`
	Reference       string               `json:"reference"`
	Status          vo.TransactionStatus `json:"status"`
	CreatedAt       time.Time            `json:"created_at"`
	CompletedAt     *time.Time           `json:"completed_at,omitempty"`
}

// NewDebitTransaction creates a new debit transaction (withdrawal)
func NewDebitTransaction(
	fromAccountID vo.AccountID,
	amount vo.Money,
	description string,
	reference string,
) (*Transaction, error) {
	if fromAccountID.IsEmpty() {
		return nil, errs.ValidationError{
			Field:   "fromAccountID",
			Message: "from account ID is required for debit transaction",
		}
	}

	if amount.IsZero() {
		return nil, errs.ErrInvalidTransactionAmount
	}

	return &Transaction{
		ID:              vo.NewTransactionID(),
		FromAccountID:   &fromAccountID,
		ToAccountID:     nil,
		TransactionType: vo.TransactionTypeDebit,
		Amount:          amount,
		Description:     strings.TrimSpace(description),
		Reference:       strings.TrimSpace(reference),
		Status:          vo.TransactionStatusPending,
		CreatedAt:       time.Now(),
	}, nil
}

// NewCreditTransaction creates a new credit transaction (deposit)
func NewCreditTransaction(
	toAccountID vo.AccountID,
	amount vo.Money,
	description string,
	reference string,
) (*Transaction, error) {
	if toAccountID.IsEmpty() {
		return nil, errs.ValidationError{
			Field:   "toAccountID",
			Message: "to account ID is required for credit transaction",
		}
	}

	if amount.IsZero() {
		return nil, errs.ErrInvalidTransactionAmount
	}

	return &Transaction{
		ID:              vo.NewTransactionID(),
		FromAccountID:   nil,
		ToAccountID:     &toAccountID,
		TransactionType: vo.TransactionTypeCredit,
		Amount:          amount,
		Description:     strings.TrimSpace(description),
		Reference:       strings.TrimSpace(reference),
		Status:          vo.TransactionStatusPending,
		CreatedAt:       time.Now(),
	}, nil
}

// NewTransferTransaction creates a new transfer transaction
func NewTransferTransaction(
	fromAccountID vo.AccountID,
	toAccountID vo.AccountID,
	amount vo.Money,
	description string,
	reference string,
) (*Transaction, error) {
	if fromAccountID.IsEmpty() {
		return nil, errs.ValidationError{
			Field:   "fromAccountID",
			Message: "from account ID is required for transfer transaction",
		}
	}

	if toAccountID.IsEmpty() {
		return nil, errs.ValidationError{
			Field:   "toAccountID",
			Message: "to account ID is required for transfer transaction",
		}
	}

	if fromAccountID.String() == toAccountID.String() {
		return nil, errs.ErrSameAccountTransfer
	}

	if amount.IsZero() {
		return nil, errs.ErrInvalidTransactionAmount
	}

	return &Transaction{
		ID:              vo.NewTransactionID(),
		FromAccountID:   &fromAccountID,
		ToAccountID:     &toAccountID,
		TransactionType: vo.TransactionTypeTransfer,
		Amount:          amount,
		Description:     strings.TrimSpace(description),
		Reference:       strings.TrimSpace(reference),
		Status:          vo.TransactionStatusPending,
		CreatedAt:       time.Now(),
	}, nil
}

// Business methods
func (t *Transaction) MarkAsCompleted() error {
	if !t.Status.CanTransitionTo(vo.TransactionStatusCompleted) {
		return errs.BusinessError{
			Code:    "INVALID_STATUS_TRANSITION",
			Message: "cannot transition from " + string(t.Status) + " to COMPLETED",
		}
	}

	now := time.Now()
	t.Status = vo.TransactionStatusCompleted
	t.CompletedAt = &now
	return nil
}

func (t *Transaction) MarkAsFailed() error {
	if !t.Status.CanTransitionTo(vo.TransactionStatusFailed) {
		return errs.BusinessError{
			Code:    "INVALID_STATUS_TRANSITION",
			Message: "cannot transition from " + string(t.Status) + " to FAILED",
		}
	}

	t.Status = vo.TransactionStatusFailed
	return nil
}

func (t *Transaction) MarkAsCancelled() error {
	if !t.Status.CanTransitionTo(vo.TransactionStatusCancelled) {
		return errs.BusinessError{
			Code:    "INVALID_STATUS_TRANSITION",
			Message: "cannot transition from " + string(t.Status) + " to CANCELLED",
		}
	}

	t.Status = vo.TransactionStatusCancelled
	return nil
}

// SetStatus sets transaction status with validation
func (t *Transaction) SetStatus(status vo.TransactionStatus) error {
	if !status.IsValid() {
		return errs.ValidationError{
			Field:   "status",
			Message: "invalid transaction status: " + string(status),
		}
	}

	if !t.Status.CanTransitionTo(status) {
		return errs.BusinessError{
			Code:    "INVALID_STATUS_TRANSITION",
			Message: "cannot transition from " + string(t.Status) + " to " + string(status),
		}
	}

	t.Status = status
	if status.IsCompleted() {
		now := time.Now()
		t.CompletedAt = &now
	}

	return nil
}
