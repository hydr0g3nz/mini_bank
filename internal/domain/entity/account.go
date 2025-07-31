package entity

import (
	"strings"
	"time"

	errs "github.com/hydr0g3nz/mini_bank/internal/domain/error"
	"github.com/hydr0g3nz/mini_bank/internal/domain/vo"
)

// Account represents a bank account
type Account struct {
	ID          vo.AccountID     `json:"id"`
	AccountName string           `json:"account_name"`
	Balance     vo.Money         `json:"balance"`
	Status      vo.AccountStatus `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// NewAccount creates a new account
func NewAccount(accountName string, initialBalance vo.Money) (*Account, error) {
	if strings.TrimSpace(accountName) == "" {
		return nil, errs.ValidationError{
			Field:   "accountName",
			Message: "account name is required",
		}
	}

	now := time.Now()
	return &Account{
		ID:          vo.NewAccountID(),
		AccountName: strings.TrimSpace(accountName),
		Balance:     initialBalance,
		Status:      vo.AccountStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Debit decreases the account balance
func (a *Account) Debit(amount vo.Money) error {
	if amount.IsZero() || !amount.IsPositive() {
		return errs.ErrInvalidTransactionAmount
	}

	newBalance, err := a.Balance.Subtract(amount)
	if err != nil {
		return err
	}

	if newBalance.Amount().IsNegative() {
		return errs.ErrInsufficientBalance
	}

	a.Balance = newBalance
	a.UpdatedAt = time.Now()
	return nil
}

// Credit increases the account balance
func (a *Account) Credit(amount vo.Money) error {
	if amount.IsZero() || !amount.IsPositive() {
		return errs.ErrInvalidTransactionAmount
	}

	newBalance, err := a.Balance.Add(amount)
	if err != nil {
		return err
	}

	a.Balance = newBalance
	a.UpdatedAt = time.Now()
	return nil
}

// Suspend suspends the account
func (a *Account) Suspend() error {
	if !a.Status.CanTransitionTo(vo.AccountStatusSuspended) {
		return errs.BusinessError{
			Code:    "INVALID_STATUS_TRANSITION",
			Message: "cannot suspend account with current status: " + string(a.Status),
		}
	}

	a.Status = vo.AccountStatusSuspended
	a.UpdatedAt = time.Now()
	return nil
}

// Activate activates the account
func (a *Account) Activate() error {
	if !a.Status.CanTransitionTo(vo.AccountStatusActive) {
		return errs.BusinessError{
			Code:    "INVALID_STATUS_TRANSITION",
			Message: "cannot activate account with current status: " + string(a.Status),
		}
	}

	a.Status = vo.AccountStatusActive
	a.UpdatedAt = time.Now()
	return nil
}

// Deactivate deactivates the account
func (a *Account) Deactivate() error {
	if !a.Status.CanTransitionTo(vo.AccountStatusInactive) {
		return errs.BusinessError{
			Code:    "INVALID_STATUS_TRANSITION",
			Message: "cannot deactivate account with current status: " + string(a.Status),
		}
	}

	a.Status = vo.AccountStatusInactive
	a.UpdatedAt = time.Now()
	return nil
}

// SetStatus sets account status with validation
func (a *Account) SetStatus(status vo.AccountStatus) error {
	if !status.IsValid() {
		return errs.ValidationError{
			Field:   "status",
			Message: "invalid account status: " + string(status),
		}
	}

	if !a.Status.CanTransitionTo(status) {
		return errs.BusinessError{
			Code:    "INVALID_STATUS_TRANSITION",
			Message: "cannot transition from " + string(a.Status) + " to " + string(status),
		}
	}

	a.Status = status
	a.UpdatedAt = time.Now()
	return nil
}

// IsActive checks if account is active
func (a *Account) IsActive() bool {
	return a.Status.IsActive()
}

// CanTransact checks if account can perform transactions
func (a *Account) CanTransact() bool {
	return a.Status.CanTransact()
}
