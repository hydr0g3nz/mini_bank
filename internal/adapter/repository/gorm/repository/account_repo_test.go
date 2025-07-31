package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

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

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Account{})
	require.NoError(t, err)

	return db
}

func createTestAccount() *entity.Account {
	money := vo.NewMoney(decimal.NewFromFloat(1000.50))
	account, _ := entity.NewAccount("Test Account", money)
	return account
}

func TestAccountRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		account *entity.Account
		wantErr bool
		errType error
	}{
		{
			name:    "successful creation",
			account: createTestAccount(),
			wantErr: false,
		},
		{
			name:    "duplicate account creation",
			account: createTestAccount(),
			wantErr: true,
			errType: errs.ErrAccountAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := repository.NewAccountRepository(db)
			ctx := context.Background()

			// For duplicate test, create the account first
			if tt.name == "duplicate account creation" {
				err := repo.Create(ctx, tt.account)
				require.NoError(t, err)
			}

			err := repo.Create(ctx, tt.account)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {

					assert.Error(t, err)
				}
			} else {
				assert.NoError(t, err)

				// Verify account was created
				var accountModel model.Account
				err = db.Where("account_id = ?", tt.account.ID.String()).First(&accountModel).Error
				assert.NoError(t, err)
				assert.Equal(t, tt.account.ID.String(), accountModel.AccountID)
				assert.Equal(t, tt.account.AccountName, accountModel.AccountName)
			}
		})
	}
}

func TestAccountRepository_GetByID(t *testing.T) {
	tests := []struct {
		name      string
		setupData bool
		accountID string
		wantErr   bool
		errType   error
	}{
		{
			name:      "successful retrieval",
			setupData: true,
			wantErr:   false,
		},
		{
			name:      "account not found",
			setupData: false,
			accountID: "2024072912345678",
			wantErr:   true,
			errType:   errs.ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := repository.NewAccountRepository(db)
			ctx := context.Background()

			var testAccount *entity.Account
			var accountID vo.AccountID

			if tt.setupData {
				testAccount = createTestAccount()
				err := repo.Create(ctx, testAccount)
				require.NoError(t, err)
				accountID = testAccount.ID
			} else {
				var err error
				accountID, err = vo.NewAccountIDFromString(tt.accountID)
				require.NoError(t, err)
			}

			account, err := repo.GetByID(ctx, accountID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, account)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, account)
				assert.Equal(t, testAccount.ID.String(), account.ID.String())
				assert.Equal(t, testAccount.AccountName, account.AccountName)
				assert.True(t, testAccount.Balance.Equal(account.Balance))
			}
		})
	}
}

func TestAccountRepository_Update(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(repo repo.AccountRepository, ctx context.Context) *entity.Account
		wantErr bool
		errType error
	}{
		{
			name: "successful update",
			setup: func(repo repo.AccountRepository, ctx context.Context) *entity.Account {
				account := createTestAccount()
				err := repo.Create(ctx, account)
				require.NoError(t, err)

				// Modify account
				account.AccountName = "Updated Account Name"
				account.UpdatedAt = time.Now()
				return account
			},
			wantErr: false,
		},
		{
			name: "update non-existent account",
			setup: func(repo repo.AccountRepository, ctx context.Context) *entity.Account {
				// Return account that doesn't exist in DB
				return createTestAccount()
			},
			wantErr: true,
			errType: errs.ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := repository.NewAccountRepository(db)
			ctx := context.Background()

			account := tt.setup(repo, ctx)
			err := repo.Update(ctx, account)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)

				// Verify update
				updatedAccount, err := repo.GetByID(ctx, account.ID)
				assert.NoError(t, err)
				assert.Equal(t, account.AccountName, updatedAccount.AccountName)
			}
		})
	}
}

func TestAccountRepository_Delete(t *testing.T) {
	tests := []struct {
		name      string
		setupData bool
		accountID string
		wantErr   bool
		errType   error
	}{
		{
			name:      "successful deletion",
			setupData: true,
			wantErr:   false,
		},
		{
			name:      "delete non-existent account",
			setupData: false,
			accountID: "2024072912345678",
			wantErr:   true,
			errType:   errs.ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := repository.NewAccountRepository(db)
			ctx := context.Background()

			var accountID vo.AccountID

			if tt.setupData {
				account := createTestAccount()
				err := repo.Create(ctx, account)
				require.NoError(t, err)
				accountID = account.ID
			} else {
				var err error
				accountID, err = vo.NewAccountIDFromString(tt.accountID)
				require.NoError(t, err)
			}

			err := repo.Delete(ctx, accountID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)

				// Verify deletion (soft delete - should not be found)
				_, err = repo.GetByID(ctx, accountID)
				assert.Error(t, err)
			}
		})
	}
}

func TestAccountRepository_List(t *testing.T) {
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
			db := setupTestDB(t)
			repo := repository.NewAccountRepository(db)
			ctx := context.Background()

			// Setup test data
			for i := 0; i < tt.setupCount; i++ {
				money := vo.NewMoney(decimal.NewFromFloat(float64(1000 + i)))
				account, err := entity.NewAccount(fmt.Sprintf("Account %d", i), money)
				require.NoError(t, err)
				err = repo.Create(ctx, account)
				require.NoError(t, err)
			}

			accounts, err := repo.List(ctx, tt.limit, tt.offset)

			assert.NoError(t, err)
			assert.Len(t, accounts, tt.wantCount)

			// Verify accounts are ordered by created_at DESC
			if len(accounts) > 1 {
				for i := 0; i < len(accounts)-1; i++ {
					assert.True(t, accounts[i].CreatedAt.After(accounts[i+1].CreatedAt) ||
						accounts[i].CreatedAt.Equal(accounts[i+1].CreatedAt))
				}
			}
		})
	}
}

func TestAccountRepository_GetByAccountName(t *testing.T) {
	tests := []struct {
		name        string
		setupData   bool
		accountName string
		wantErr     bool
		errType     error
	}{
		{
			name:        "successful retrieval by name",
			setupData:   true,
			accountName: "Test Account",
			wantErr:     false,
		},
		{
			name:        "account not found by name",
			setupData:   false,
			accountName: "Non-existent Account",
			wantErr:     true,
			errType:     errs.ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := repository.NewAccountRepository(db)
			ctx := context.Background()

			var testAccount *entity.Account

			if tt.setupData {
				testAccount = createTestAccount()
				err := repo.Create(ctx, testAccount)
				require.NoError(t, err)
			}

			account, err := repo.GetByAccountName(ctx, tt.accountName)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, account)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, account)
				assert.Equal(t, tt.accountName, account.AccountName)
				assert.Equal(t, testAccount.ID.String(), account.ID.String())
			}
		})
	}
}
