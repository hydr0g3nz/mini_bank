package entity

import (
	"testing"
	"time"

	errs "github.com/hydr0g3nz/mini_bank/internal/domain/error"
	"github.com/hydr0g3nz/mini_bank/internal/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDebitTransaction(t *testing.T) {
	validAccountID := vo.NewAccountID()
	amount := vo.NewMoneyFromFloat(100.0)

	tests := []struct {
		name          string
		fromAccountID vo.AccountID
		amount        vo.Money
		description   string
		reference     string
		expectError   bool
		errorType     interface{}
	}{
		{
			name:          "Valid debit transaction",
			fromAccountID: validAccountID,
			amount:        amount,
			description:   "Test debit",
			reference:     "REF001",
			expectError:   false,
		},
		{
			name:          "Empty from account ID",
			fromAccountID: vo.AccountID{},
			amount:        amount,
			description:   "Test debit",
			reference:     "REF001",
			expectError:   true,
			errorType:     errs.ValidationError{},
		},
		{
			name:          "Zero amount",
			fromAccountID: validAccountID,
			amount:        vo.ZeroMoney(),
			description:   "Test debit",
			reference:     "REF001",
			expectError:   true,
			errorType:     errs.ErrInvalidTransactionAmount,
		},
		{
			name:          "With whitespace in fields",
			fromAccountID: validAccountID,
			amount:        amount,
			description:   "  Test debit  ",
			reference:     "  REF001  ",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction, err := NewDebitTransaction(tt.fromAccountID, tt.amount, tt.description, tt.reference)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorType != nil {
					assert.IsType(t, tt.errorType, err)
				}
				assert.Nil(t, transaction)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, transaction)
				assert.NotEmpty(t, transaction.ID.String())
				assert.Equal(t, &tt.fromAccountID, transaction.FromAccountID)
				assert.Nil(t, transaction.ToAccountID)
				assert.Equal(t, vo.TransactionTypeDebit, transaction.TransactionType)
				assert.True(t, transaction.Amount.Equal(tt.amount))
				assert.Equal(t, "Test debit", transaction.Description)
				assert.Equal(t, "REF001", transaction.Reference)
				assert.Equal(t, vo.TransactionStatusPending, transaction.Status)
				assert.WithinDuration(t, time.Now(), transaction.CreatedAt, time.Second)
				assert.Nil(t, transaction.CompletedAt)
			}
		})
	}
}

func TestNewCreditTransaction(t *testing.T) {
	validAccountID := vo.NewAccountID()
	amount := vo.NewMoneyFromFloat(100.0)

	tests := []struct {
		name        string
		toAccountID vo.AccountID
		amount      vo.Money
		description string
		reference   string
		expectError bool
		errorType   interface{}
	}{
		{
			name:        "Valid credit transaction",
			toAccountID: validAccountID,
			amount:      amount,
			description: "Test credit",
			reference:   "REF002",
			expectError: false,
		},
		{
			name:        "Empty to account ID",
			toAccountID: vo.AccountID{},
			amount:      amount,
			description: "Test credit",
			reference:   "REF002",
			expectError: true,
			errorType:   errs.ValidationError{},
		},
		{
			name:        "Zero amount",
			toAccountID: validAccountID,
			amount:      vo.ZeroMoney(),
			description: "Test credit",
			reference:   "REF002",
			expectError: true,
			errorType:   errs.ErrInvalidTransactionAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction, err := NewCreditTransaction(tt.toAccountID, tt.amount, tt.description, tt.reference)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorType != nil {
					assert.IsType(t, tt.errorType, err)
				}
				assert.Nil(t, transaction)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, transaction)
				assert.NotEmpty(t, transaction.ID.String())
				assert.Nil(t, transaction.FromAccountID)
				assert.Equal(t, &tt.toAccountID, transaction.ToAccountID)
				assert.Equal(t, vo.TransactionTypeCredit, transaction.TransactionType)
				assert.True(t, transaction.Amount.Equal(tt.amount))
				assert.Equal(t, "Test credit", transaction.Description)
				assert.Equal(t, "REF002", transaction.Reference)
				assert.Equal(t, vo.TransactionStatusPending, transaction.Status)
				assert.WithinDuration(t, time.Now(), transaction.CreatedAt, time.Second)
				assert.Nil(t, transaction.CompletedAt)
			}
		})
	}
}

func TestNewTransferTransaction(t *testing.T) {
	fromAccountID := vo.NewAccountID()
	toAccountID := vo.NewAccountID()
	amount := vo.NewMoneyFromFloat(100.0)

	tests := []struct {
		name          string
		fromAccountID vo.AccountID
		toAccountID   vo.AccountID
		amount        vo.Money
		description   string
		reference     string
		expectError   bool
		errorType     interface{}
	}{
		{
			name:          "Valid transfer transaction",
			fromAccountID: fromAccountID,
			toAccountID:   toAccountID,
			amount:        amount,
			description:   "Test transfer",
			reference:     "REF003",
			expectError:   false,
		},
		{
			name:          "Empty from account ID",
			fromAccountID: vo.AccountID{},
			toAccountID:   toAccountID,
			amount:        amount,
			description:   "Test transfer",
			reference:     "REF003",
			expectError:   true,
			errorType:     errs.ValidationError{},
		},
		{
			name:          "Empty to account ID",
			fromAccountID: fromAccountID,
			toAccountID:   vo.AccountID{},
			amount:        amount,
			description:   "Test transfer",
			reference:     "REF003",
			expectError:   true,
			errorType:     errs.ValidationError{},
		},
		{
			name:          "Same from and to account",
			fromAccountID: fromAccountID,
			toAccountID:   fromAccountID,
			amount:        amount,
			description:   "Test transfer",
			reference:     "REF003",
			expectError:   true,
			errorType:     errs.ErrSameAccountTransfer,
		},
		{
			name:          "Zero amount",
			fromAccountID: fromAccountID,
			toAccountID:   toAccountID,
			amount:        vo.ZeroMoney(),
			description:   "Test transfer",
			reference:     "REF003",
			expectError:   true,
			errorType:     errs.ErrInvalidTransactionAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction, err := NewTransferTransaction(tt.fromAccountID, tt.toAccountID, tt.amount, tt.description, tt.reference)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorType != nil {
					assert.IsType(t, tt.errorType, err)
				}
				assert.Nil(t, transaction)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, transaction)
				assert.NotEmpty(t, transaction.ID.String())
				assert.Equal(t, &tt.fromAccountID, transaction.FromAccountID)
				assert.Equal(t, &tt.toAccountID, transaction.ToAccountID)
				assert.Equal(t, vo.TransactionTypeTransfer, transaction.TransactionType)
				assert.True(t, transaction.Amount.Equal(tt.amount))
				assert.Equal(t, "Test transfer", transaction.Description)
				assert.Equal(t, "REF003", transaction.Reference)
				assert.Equal(t, vo.TransactionStatusPending, transaction.Status)
				assert.WithinDuration(t, time.Now(), transaction.CreatedAt, time.Second)
				assert.Nil(t, transaction.CompletedAt)
			}
		})
	}
}

func TestTransaction_MarkAsCompleted(t *testing.T) {
	fromAccountID := vo.NewAccountID()
	amount := vo.NewMoneyFromFloat(100.0)

	tests := []struct {
		name          string
		initialStatus vo.TransactionStatus
		expectError   bool
	}{
		{
			name:          "Mark pending transaction as completed",
			initialStatus: vo.TransactionStatusPending,
			expectError:   false,
		},
		{
			name:          "Mark completed transaction as completed",
			initialStatus: vo.TransactionStatusCompleted,
			expectError:   true,
		},
		{
			name:          "Mark failed transaction as completed",
			initialStatus: vo.TransactionStatusFailed,
			expectError:   true,
		},
		{
			name:          "Mark cancelled transaction as completed",
			initialStatus: vo.TransactionStatusCancelled,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction, err := NewDebitTransaction(fromAccountID, amount, "Test", "REF")
			require.NoError(t, err)

			transaction.Status = tt.initialStatus

			err = transaction.MarkAsCompleted()

			if tt.expectError {
				require.Error(t, err)
				assert.IsType(t, errs.BusinessError{}, err)
				assert.Equal(t, tt.initialStatus, transaction.Status)
				assert.Nil(t, transaction.CompletedAt)
			} else {
				require.NoError(t, err)
				assert.Equal(t, vo.TransactionStatusCompleted, transaction.Status)
				assert.NotNil(t, transaction.CompletedAt)
				assert.WithinDuration(t, time.Now(), *transaction.CompletedAt, time.Second)
			}
		})
	}
}

func TestTransaction_MarkAsFailed(t *testing.T) {
	fromAccountID := vo.NewAccountID()
	amount := vo.NewMoneyFromFloat(100.0)

	tests := []struct {
		name          string
		initialStatus vo.TransactionStatus
		expectError   bool
	}{
		{
			name:          "Mark pending transaction as failed",
			initialStatus: vo.TransactionStatusPending,
			expectError:   false,
		},
		{
			name:          "Mark completed transaction as failed",
			initialStatus: vo.TransactionStatusCompleted,
			expectError:   true,
		},
		{
			name:          "Mark cancelled transaction as failed",
			initialStatus: vo.TransactionStatusCancelled,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction, err := NewDebitTransaction(fromAccountID, amount, "Test", "REF")
			require.NoError(t, err)

			transaction.Status = tt.initialStatus

			err = transaction.MarkAsFailed()

			if tt.expectError {
				require.Error(t, err)
				assert.IsType(t, errs.BusinessError{}, err)
				assert.Equal(t, tt.initialStatus, transaction.Status)
			} else {
				require.NoError(t, err)
				assert.Equal(t, vo.TransactionStatusFailed, transaction.Status)
			}
		})
	}
}

func TestTransaction_MarkAsCancelled(t *testing.T) {
	fromAccountID := vo.NewAccountID()
	amount := vo.NewMoneyFromFloat(100.0)

	tests := []struct {
		name          string
		initialStatus vo.TransactionStatus
		expectError   bool
	}{
		{
			name:          "Mark pending transaction as cancelled",
			initialStatus: vo.TransactionStatusPending,
			expectError:   false,
		},
		{
			name:          "Mark failed transaction as cancelled",
			initialStatus: vo.TransactionStatusFailed,
			expectError:   false,
		},
		{
			name:          "Mark completed transaction as cancelled",
			initialStatus: vo.TransactionStatusCompleted,
			expectError:   true,
		},
		{
			name:          "Mark cancelled transaction as cancelled",
			initialStatus: vo.TransactionStatusCancelled,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction, err := NewDebitTransaction(fromAccountID, amount, "Test", "REF")
			require.NoError(t, err)

			transaction.Status = tt.initialStatus

			err = transaction.MarkAsCancelled()

			if tt.expectError {
				require.Error(t, err)
				assert.IsType(t, errs.BusinessError{}, err)
				assert.Equal(t, tt.initialStatus, transaction.Status)
			} else {
				require.NoError(t, err)
				assert.Equal(t, vo.TransactionStatusCancelled, transaction.Status)
			}
		})
	}
}

func TestTransaction_SetStatus(t *testing.T) {
	fromAccountID := vo.NewAccountID()
	amount := vo.NewMoneyFromFloat(100.0)

	tests := []struct {
		name          string
		initialStatus vo.TransactionStatus
		targetStatus  vo.TransactionStatus
		expectError   bool
		errorType     interface{}
	}{
		{
			name:          "Valid status transition",
			initialStatus: vo.TransactionStatusPending,
			targetStatus:  vo.TransactionStatusCompleted,
			expectError:   false,
		},
		{
			name:          "Invalid status",
			initialStatus: vo.TransactionStatusPending,
			targetStatus:  vo.TransactionStatus("INVALID"),
			expectError:   true,
			errorType:     errs.ValidationError{},
		},
		{
			name:          "Invalid transition",
			initialStatus: vo.TransactionStatusCompleted,
			targetStatus:  vo.TransactionStatusPending,
			expectError:   true,
			errorType:     errs.BusinessError{},
		},
		{
			name:          "Set completed status with timestamp",
			initialStatus: vo.TransactionStatusPending,
			targetStatus:  vo.TransactionStatusCompleted,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction, err := NewDebitTransaction(fromAccountID, amount, "Test", "REF")
			require.NoError(t, err)

			transaction.Status = tt.initialStatus

			err = transaction.SetStatus(tt.targetStatus)

			if tt.expectError {
				require.Error(t, err)
				assert.IsType(t, tt.errorType, err)
				assert.Equal(t, tt.initialStatus, transaction.Status)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.targetStatus, transaction.Status)

				if tt.targetStatus.IsCompleted() {
					assert.NotNil(t, transaction.CompletedAt)
					assert.WithinDuration(t, time.Now(), *transaction.CompletedAt, time.Second)
				}
			}
		})
	}
}
