package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hydr0g3nz/mini_bank/internal/adapter/repository/gorm/model"
	"github.com/hydr0g3nz/mini_bank/internal/adapter/repository/gorm/repository"
	"github.com/hydr0g3nz/mini_bank/internal/domain/entity"
	errs "github.com/hydr0g3nz/mini_bank/internal/domain/error"
	repo "github.com/hydr0g3nz/mini_bank/internal/domain/repository"
	"github.com/hydr0g3nz/mini_bank/internal/domain/vo"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTransactionTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Account{}, &model.Transaction{})
	require.NoError(t, err)

	return db
}

func createTestTransactions() (*entity.Transaction, *entity.Transaction, *entity.Transaction) {
	fromAccountID := vo.NewAccountID()
	toAccountID := vo.NewAccountID()
	amount := vo.NewMoney(decimal.NewFromFloat(100.50))

	debitTxn, _ := entity.NewDebitTransaction(fromAccountID, amount, "Test debit", "REF001")
	creditTxn, _ := entity.NewCreditTransaction(toAccountID, amount, "Test credit", "REF002")
	transferTxn, _ := entity.NewTransferTransaction(fromAccountID, toAccountID, amount, "Test transfer", "REF003")

	return debitTxn, creditTxn, transferTxn
}

func TestTransactionRepository_Create(t *testing.T) {
	tests := []struct {
		name        string
		transaction *entity.Transaction
		wantErr     bool
	}{
		{
			name: "successful debit transaction creation",
			transaction: func() *entity.Transaction {
				debit, _, _ := createTestTransactions()
				return debit
			}(),
			wantErr: false,
		},
		{
			name: "successful credit transaction creation",
			transaction: func() *entity.Transaction {
				_, credit, _ := createTestTransactions()
				return credit
			}(),
			wantErr: false,
		},
		{
			name: "successful transfer transaction creation",
			transaction: func() *entity.Transaction {
				_, _, transfer := createTestTransactions()
				return transfer
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTransactionTestDB(t)
			repo := repository.NewTransactionRepository(db)
			ctx := context.Background()

			err := repo.Create(ctx, tt.transaction)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify transaction was created
				var txnModel model.Transaction
				err = db.Where("transaction_id = ?", tt.transaction.ID.String()).First(&txnModel).Error
				assert.NoError(t, err)
				assert.Equal(t, tt.transaction.ID.String(), txnModel.TransactionID)
				assert.Equal(t, string(tt.transaction.TransactionType), txnModel.TransactionType)
			}
		})
	}
}

func TestTransactionRepository_GetByID(t *testing.T) {
	tests := []struct {
		name          string
		setupData     bool
		transactionID string
		wantErr       bool
		errType       error
	}{
		{
			name:      "successful retrieval",
			setupData: true,
			wantErr:   false,
		},
		{
			name:          "transaction not found",
			setupData:     false,
			transactionID: "TXN20240729143045123456",
			wantErr:       true,
			errType:       errs.ErrTransactionNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTransactionTestDB(t)
			repo := repository.NewTransactionRepository(db)
			ctx := context.Background()

			var testTransaction *entity.Transaction
			var transactionID vo.TransactionID

			if tt.setupData {
				testTransaction, _, _ = createTestTransactions()
				err := repo.Create(ctx, testTransaction)
				require.NoError(t, err)
				transactionID = testTransaction.ID
			} else {
				var err error
				transactionID, err = vo.NewTransactionIDFromString(tt.transactionID)
				require.NoError(t, err)
			}

			transaction, err := repo.GetByID(ctx, transactionID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, transaction)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, transaction)
				assert.Equal(t, testTransaction.ID.String(), transaction.ID.String())
				assert.Equal(t, testTransaction.TransactionType, transaction.TransactionType)
				assert.True(t, testTransaction.Amount.Equal(transaction.Amount))
			}
		})
	}
}

func TestTransactionRepository_Update(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(repo repo.TransactionRepository, ctx context.Context) *entity.Transaction
		wantErr bool
		errType error
	}{
		{
			name: "successful update",
			setup: func(repo repo.TransactionRepository, ctx context.Context) *entity.Transaction {
				transaction, _, _ := createTestTransactions()
				err := repo.Create(ctx, transaction)
				require.NoError(t, err)

				// Mark as completed
				err = transaction.MarkAsCompleted()
				require.NoError(t, err)

				return transaction
			},
			wantErr: false,
		},
		{
			name: "update non-existent transaction",
			setup: func(repo repo.TransactionRepository, ctx context.Context) *entity.Transaction {
				// Return transaction that doesn't exist in DB
				transaction, _, _ := createTestTransactions()
				return transaction
			},
			wantErr: true,
			errType: errs.ErrTransactionNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTransactionTestDB(t)
			repo := repository.NewTransactionRepository(db)
			ctx := context.Background()

			transaction := tt.setup(repo, ctx)
			err := repo.Update(ctx, transaction)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)

				// Verify update
				updatedTransaction, err := repo.GetByID(ctx, transaction.ID)
				assert.NoError(t, err)
				assert.Equal(t, transaction.Status, updatedTransaction.Status)
				if transaction.CompletedAt != nil {
					assert.NotNil(t, updatedTransaction.CompletedAt)
				}
			}
		})
	}
}

func TestTransactionRepository_List(t *testing.T) {
	tests := []struct {
		name       string
		setupCount int
		limit      int
		offset     int
		wantCount  int
	}{
		{
			name:       "list with limit",
			setupCount: 5,
			limit:      3,
			offset:     0,
			wantCount:  3,
		},
		{
			name:       "list with offset",
			setupCount: 5,
			limit:      10,
			offset:     2,
			wantCount:  3,
		},
		{
			name:       "empty list",
			setupCount: 0,
			limit:      10,
			offset:     0,
			wantCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTransactionTestDB(t)
			repo := repository.NewTransactionRepository(db)
			ctx := context.Background()

			// Setup test data
			for i := 0; i < tt.setupCount; i++ {
				fromAccountID := vo.NewAccountID()
				amount := vo.NewMoney(decimal.NewFromFloat(float64(100 + i)))
				transaction, err := entity.NewDebitTransaction(
					fromAccountID,
					amount,
					fmt.Sprintf("Test transaction %d", i),
					fmt.Sprintf("REF%03d", i),
				)
				require.NoError(t, err)
				err = repo.Create(ctx, transaction)
				require.NoError(t, err)
			}

			transactions, err := repo.List(ctx, tt.limit, tt.offset)

			assert.NoError(t, err)
			assert.Len(t, transactions, tt.wantCount)

			// Verify transactions are ordered by created_at DESC
			if len(transactions) > 1 {
				for i := 0; i < len(transactions)-1; i++ {
					assert.True(t, transactions[i].CreatedAt.After(transactions[i+1].CreatedAt) ||
						transactions[i].CreatedAt.Equal(transactions[i+1].CreatedAt))
				}
			}
		})
	}
}

func TestTransactionRepository_GetByAccountID(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(repo repo.TransactionRepository, ctx context.Context) vo.AccountID
		limit     int
		offset    int
		wantCount int
	}{
		{
			name: "get transactions for account (as from_account)",
			setup: func(repo repo.TransactionRepository, ctx context.Context) vo.AccountID {
				fromAccountID := vo.NewAccountID()
				toAccountID := vo.NewAccountID()
				amount := vo.NewMoney(decimal.NewFromFloat(100))

				// Create multiple transactions for the same account
				for i := 0; i < 3; i++ {
					transaction, err := entity.NewTransferTransaction(
						fromAccountID,
						toAccountID,
						amount,
						fmt.Sprintf("Transfer %d", i),
						fmt.Sprintf("REF%d", i),
					)
					require.NoError(t, err)
					err = repo.Create(ctx, transaction)
					require.NoError(t, err)
				}

				return fromAccountID
			},
			limit:     10,
			offset:    0,
			wantCount: 3,
		},
		{
			name: "get transactions for account (as to_account)",
			setup: func(repo repo.TransactionRepository, ctx context.Context) vo.AccountID {
				fromAccountID := vo.NewAccountID()
				toAccountID := vo.NewAccountID()
				amount := vo.NewMoney(decimal.NewFromFloat(100))

				// Create multiple transactions for the same account
				for i := 0; i < 2; i++ {
					transaction, err := entity.NewTransferTransaction(
						fromAccountID,
						toAccountID,
						amount,
						fmt.Sprintf("Transfer %d", i),
						fmt.Sprintf("REF%d", i),
					)
					require.NoError(t, err)
					err = repo.Create(ctx, transaction)
					require.NoError(t, err)
				}

				return toAccountID
			},
			limit:     10,
			offset:    0,
			wantCount: 2,
		},
		{
			name: "no transactions for account",
			setup: func(repo repo.TransactionRepository, ctx context.Context) vo.AccountID {
				return vo.NewAccountID() // Return a new account ID that has no transactions
			},
			limit:     10,
			offset:    0,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTransactionTestDB(t)
			repo := repository.NewTransactionRepository(db)
			ctx := context.Background()

			accountID := tt.setup(repo, ctx)
			transactions, err := repo.GetByAccountID(ctx, accountID, tt.limit, tt.offset)

			assert.NoError(t, err)
			assert.Len(t, transactions, tt.wantCount)

			// Verify all transactions involve the specified account
			for _, txn := range transactions {
				accountInvolved := false
				if txn.FromAccountID != nil && txn.FromAccountID.String() == accountID.String() {
					accountInvolved = true
				}
				if txn.ToAccountID != nil && txn.ToAccountID.String() == accountID.String() {
					accountInvolved = true
				}
				assert.True(t, accountInvolved, "Transaction should involve the specified account")
			}
		})
	}
}

func TestTransactionRepository_GetByStatus(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(repo repo.TransactionRepository, ctx context.Context)
		status    vo.TransactionStatus
		limit     int
		offset    int
		wantCount int
	}{
		{
			name: "get pending transactions",
			setup: func(repo repo.TransactionRepository, ctx context.Context) {
				fromAccountID := vo.NewAccountID()
				amount := vo.NewMoney(decimal.NewFromFloat(100))

				// Create pending transactions
				for i := 0; i < 3; i++ {
					transaction, err := entity.NewDebitTransaction(
						fromAccountID,
						amount,
						fmt.Sprintf("Debit %d", i),
						fmt.Sprintf("REF%d", i),
					)
					require.NoError(t, err)
					err = repo.Create(ctx, transaction)
					require.NoError(t, err)
				}

				// Create one completed transaction
				transaction, err := entity.NewDebitTransaction(
					fromAccountID,
					amount,
					"Completed debit",
					"REF999",
				)
				require.NoError(t, err)
				err = transaction.MarkAsCompleted()
				require.NoError(t, err)
				err = repo.Create(ctx, transaction)
				require.NoError(t, err)
			},
			status:    vo.TransactionStatusPending,
			limit:     10,
			offset:    0,
			wantCount: 3,
		},
		{
			name: "get completed transactions",
			setup: func(repo repo.TransactionRepository, ctx context.Context) {
				fromAccountID := vo.NewAccountID()
				amount := vo.NewMoney(decimal.NewFromFloat(100))

				// Create and complete transactions
				for i := 0; i < 2; i++ {
					transaction, err := entity.NewDebitTransaction(
						fromAccountID,
						amount,
						fmt.Sprintf("Debit %d", i),
						fmt.Sprintf("REF%d", i),
					)
					require.NoError(t, err)
					err = transaction.MarkAsCompleted()
					require.NoError(t, err)
					err = repo.Create(ctx, transaction)
					require.NoError(t, err)
				}
			},
			status:    vo.TransactionStatusCompleted,
			limit:     10,
			offset:    0,
			wantCount: 2,
		},
		{
			name: "no transactions with status",
			setup: func(repo repo.TransactionRepository, ctx context.Context) {
				// Create only pending transactions
				fromAccountID := vo.NewAccountID()
				amount := vo.NewMoney(decimal.NewFromFloat(100))
				transaction, err := entity.NewDebitTransaction(fromAccountID, amount, "Debit", "REF001")
				require.NoError(t, err)
				err = repo.Create(ctx, transaction)
				require.NoError(t, err)
			},
			status:    vo.TransactionStatusFailed,
			limit:     10,
			offset:    0,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTransactionTestDB(t)
			repo := repository.NewTransactionRepository(db)
			ctx := context.Background()

			tt.setup(repo, ctx)
			transactions, err := repo.GetByStatus(ctx, tt.status, tt.limit, tt.offset)

			assert.NoError(t, err)
			assert.Len(t, transactions, tt.wantCount)

			// Verify all transactions have the specified status
			for _, txn := range transactions {
				assert.Equal(t, tt.status, txn.Status)
			}

			// Verify transactions are ordered by created_at DESC
			if len(transactions) > 1 {
				for i := 0; i < len(transactions)-1; i++ {
					assert.True(t, transactions[i].CreatedAt.After(transactions[i+1].CreatedAt) ||
						transactions[i].CreatedAt.Equal(transactions[i+1].CreatedAt))
				}
			}
		})
	}
}
