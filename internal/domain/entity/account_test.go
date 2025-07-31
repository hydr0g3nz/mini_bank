package entity

import (
	"testing"
	"time"

	errs "github.com/hydr0g3nz/mini_bank/internal/domain/error"
	"github.com/hydr0g3nz/mini_bank/internal/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccount(t *testing.T) {
	tests := []struct {
		name           string
		accountName    string
		initialBalance vo.Money
		expectError    bool
		errorType      interface{}
	}{
		{
			name:           "Valid account creation",
			accountName:    "Test Account",
			initialBalance: vo.NewMoneyFromFloat(100.0),
			expectError:    false,
		},
		{
			name:           "Account with zero balance",
			accountName:    "Zero Balance Account",
			initialBalance: vo.ZeroMoney(),
			expectError:    false,
		},
		{
			name:           "Empty account name",
			accountName:    "",
			initialBalance: vo.NewMoneyFromFloat(100.0),
			expectError:    true,
			errorType:      errs.ValidationError{},
		},
		{
			name:           "Whitespace only account name",
			accountName:    "   ",
			initialBalance: vo.NewMoneyFromFloat(100.0),
			expectError:    true,
			errorType:      errs.ValidationError{},
		},
		{
			name:           "Account with trimmed name",
			accountName:    "  Test Account  ",
			initialBalance: vo.NewMoneyFromFloat(50.0),
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := NewAccount(tt.accountName, tt.initialBalance)

			if tt.expectError {
				require.Error(t, err)
				assert.IsType(t, tt.errorType, err)
				assert.Nil(t, account)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, account)
				assert.NotEmpty(t, account.ID.String())
				assert.NotEmpty(t, account.AccountName)
				assert.True(t, account.Balance.Equal(tt.initialBalance))
				assert.Equal(t, vo.AccountStatusActive, account.Status)
				assert.WithinDuration(t, time.Now(), account.CreatedAt, time.Second)
				assert.WithinDuration(t, time.Now(), account.UpdatedAt, time.Second)
			}
		})
	}
}

func TestAccount_Debit(t *testing.T) {
	tests := []struct {
		name            string
		initialBalance  vo.Money
		debitAmount     vo.Money
		expectError     bool
		expectedBalance vo.Money
		errorType       error
	}{
		{
			name:            "Valid debit",
			initialBalance:  vo.NewMoneyFromFloat(100.0),
			debitAmount:     vo.NewMoneyFromFloat(50.0),
			expectError:     false,
			expectedBalance: vo.NewMoneyFromFloat(50.0),
		},
		{
			name:            "Debit full balance",
			initialBalance:  vo.NewMoneyFromFloat(100.0),
			debitAmount:     vo.NewMoneyFromFloat(100.0),
			expectError:     false,
			expectedBalance: vo.ZeroMoney(),
		},
		{
			name:           "Insufficient balance",
			initialBalance: vo.NewMoneyFromFloat(50.0),
			debitAmount:    vo.NewMoneyFromFloat(100.0),
			expectError:    true,
			errorType:      errs.ErrInsufficientBalance,
		},
		{
			name:           "Zero debit amount",
			initialBalance: vo.NewMoneyFromFloat(100.0),
			debitAmount:    vo.ZeroMoney(),
			expectError:    true,
			errorType:      errs.ErrInvalidTransactionAmount,
		},
		{
			name:           "Negative debit amount",
			initialBalance: vo.NewMoneyFromFloat(100.0),
			debitAmount:    vo.NewMoneyFromFloat(-50.0),
			expectError:    true,
			errorType:      errs.ErrInvalidTransactionAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := NewAccount("Test Account", tt.initialBalance)
			require.NoError(t, err)

			originalUpdatedAt := account.UpdatedAt
			time.Sleep(time.Millisecond * 10) // Ensure time difference

			err = account.Debit(tt.debitAmount)

			if tt.expectError {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.errorType)
				assert.True(t, account.Balance.Equal(tt.initialBalance)) // Balance unchanged
				assert.Equal(t, originalUpdatedAt, account.UpdatedAt)    // UpdatedAt unchanged
			} else {
				require.NoError(t, err)
				assert.True(t, account.Balance.Equal(tt.expectedBalance))
				assert.True(t, account.UpdatedAt.After(originalUpdatedAt))
			}
		})
	}
}

func TestAccount_Credit(t *testing.T) {
	tests := []struct {
		name            string
		initialBalance  vo.Money
		creditAmount    vo.Money
		expectError     bool
		expectedBalance vo.Money
		errorType       error
	}{
		{
			name:            "Valid credit",
			initialBalance:  vo.NewMoneyFromFloat(100.0),
			creditAmount:    vo.NewMoneyFromFloat(50.0),
			expectError:     false,
			expectedBalance: vo.NewMoneyFromFloat(150.0),
		},
		{
			name:            "Credit to zero balance",
			initialBalance:  vo.ZeroMoney(),
			creditAmount:    vo.NewMoneyFromFloat(100.0),
			expectError:     false,
			expectedBalance: vo.NewMoneyFromFloat(100.0),
		},
		{
			name:           "Zero credit amount",
			initialBalance: vo.NewMoneyFromFloat(100.0),
			creditAmount:   vo.ZeroMoney(),
			expectError:    true,
			errorType:      errs.ErrInvalidTransactionAmount,
		},
		{
			name:           "Negative credit amount",
			initialBalance: vo.NewMoneyFromFloat(100.0),
			creditAmount:   vo.NewMoneyFromFloat(-50.0),
			expectError:    true,
			errorType:      errs.ErrInvalidTransactionAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := NewAccount("Test Account", tt.initialBalance)
			require.NoError(t, err)

			originalUpdatedAt := account.UpdatedAt
			time.Sleep(time.Millisecond * 10) // Ensure time difference

			err = account.Credit(tt.creditAmount)

			if tt.expectError {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.errorType)
				assert.True(t, account.Balance.Equal(tt.initialBalance)) // Balance unchanged
				assert.Equal(t, originalUpdatedAt, account.UpdatedAt)    // UpdatedAt unchanged
			} else {
				require.NoError(t, err)
				assert.True(t, account.Balance.Equal(tt.expectedBalance))
				assert.True(t, account.UpdatedAt.After(originalUpdatedAt))
			}
		})
	}
}

func TestAccount_StatusTransitions(t *testing.T) {
	account, err := NewAccount("Test Account", vo.NewMoneyFromFloat(100.0))
	require.NoError(t, err)

	t.Run("Suspend active account", func(t *testing.T) {
		err := account.Suspend()
		require.NoError(t, err)
		assert.Equal(t, vo.AccountStatusSuspended, account.Status)
	})

	t.Run("Activate suspended account", func(t *testing.T) {
		err := account.Activate()
		require.NoError(t, err)
		assert.Equal(t, vo.AccountStatusActive, account.Status)
	})

	t.Run("Deactivate active account", func(t *testing.T) {
		err := account.Deactivate()
		require.NoError(t, err)
		assert.Equal(t, vo.AccountStatusInactive, account.Status)
	})

	t.Run("Activate inactive account", func(t *testing.T) {
		err := account.Activate()
		require.NoError(t, err)
		assert.Equal(t, vo.AccountStatusActive, account.Status)
	})
}

func TestAccount_StatusTransitionErrors(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus vo.AccountStatus
		operation     func(*Account) error
		expectError   bool
	}{
		{
			name:          "Invalid status transition",
			initialStatus: vo.AccountStatusSuspended,
			operation: func(a *Account) error {
				// Try to suspend already suspended account
				return a.Suspend()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := NewAccount("Test Account", vo.NewMoneyFromFloat(100.0))
			require.NoError(t, err)

			// Set initial status
			account.Status = tt.initialStatus

			err = tt.operation(account)

			if tt.expectError {
				require.Error(t, err)
				assert.IsType(t, errs.BusinessError{}, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAccount_SetStatus(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus vo.AccountStatus
		targetStatus  vo.AccountStatus
		expectError   bool
		errorType     interface{}
	}{
		{
			name:          "Valid status transition",
			initialStatus: vo.AccountStatusActive,
			targetStatus:  vo.AccountStatusSuspended,
			expectError:   false,
		},
		{
			name:          "Invalid status",
			initialStatus: vo.AccountStatusActive,
			targetStatus:  vo.AccountStatus("INVALID"),
			expectError:   true,
			errorType:     errs.ValidationError{},
		},
		{
			name:          "Invalid transition",
			initialStatus: vo.AccountStatusSuspended,
			targetStatus:  vo.AccountStatusSuspended,
			expectError:   true,
			errorType:     errs.BusinessError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := NewAccount("Test Account", vo.NewMoneyFromFloat(100.0))
			require.NoError(t, err)

			account.Status = tt.initialStatus
			originalUpdatedAt := account.UpdatedAt
			time.Sleep(time.Millisecond * 10)

			err = account.SetStatus(tt.targetStatus)

			if tt.expectError {
				require.Error(t, err)
				assert.IsType(t, tt.errorType, err)
				assert.Equal(t, tt.initialStatus, account.Status)
				assert.Equal(t, originalUpdatedAt, account.UpdatedAt)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.targetStatus, account.Status)
				assert.True(t, account.UpdatedAt.After(originalUpdatedAt))
			}
		})
	}
}

func TestAccount_StatusChecks(t *testing.T) {
	account, err := NewAccount("Test Account", vo.NewMoneyFromFloat(100.0))
	require.NoError(t, err)

	// Test active account
	assert.True(t, account.IsActive())
	assert.True(t, account.CanTransact())

	// Test suspended account
	account.Status = vo.AccountStatusSuspended
	assert.False(t, account.IsActive())
	assert.False(t, account.CanTransact())

	// Test inactive account
	account.Status = vo.AccountStatusInactive
	assert.False(t, account.IsActive())
	assert.False(t, account.CanTransact())
}
