package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hydr0g3nz/mini_bank/internal/application/dto"
	"github.com/hydr0g3nz/mini_bank/internal/domain/entity"
	errs "github.com/hydr0g3nz/mini_bank/internal/domain/error"
	"github.com/hydr0g3nz/mini_bank/internal/domain/infra"
	"github.com/hydr0g3nz/mini_bank/internal/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock structs
type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) Create(ctx context.Context, account *entity.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAccountRepository) GetByID(ctx context.Context, id vo.AccountID) (*entity.Account, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func (m *MockAccountRepository) Update(ctx context.Context, account *entity.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAccountRepository) Delete(ctx context.Context, id vo.AccountID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAccountRepository) List(ctx context.Context, limit, offset int) ([]*entity.Account, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*entity.Account), args.Error(1)
}

func (m *MockAccountRepository) GetByAccountName(ctx context.Context, accountName string) (*entity.Account, error) {
	args := m.Called(ctx, accountName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Account), args.Error(1)
}

type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheService) Get(ctx context.Context, key string, dest interface{}) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func (m *MockCacheService) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Debugf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Info(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Warn(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Warnf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Error(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Fatal(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Fatalf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) With(fields ...interface{}) infra.Logger {
	args := m.Called(fields)
	return args.Get(0).(infra.Logger)
}

func (m *MockLogger) Sync() error {
	args := m.Called()
	return args.Error(0)
}

// Test fixtures
func createTestAccount() *entity.Account {
	account, _ := entity.NewAccount("Test Account", vo.NewMoneyFromFloat(1000.0))
	return account
}

func TestAccountUseCase_CreateAccount(t *testing.T) {
	tests := []struct {
		name           string
		request        dto.CreateAccountRequest
		setupMocks     func(*MockAccountRepository, *MockCacheService, *MockLogger)
		expectedError  error
		validateResult func(*testing.T, *dto.AccountResponse)
	}{
		{
			name: "success_create_account",
			request: dto.CreateAccountRequest{
				AccountName:    "Test Account",
				InitialBalance: 1000.0,
			},
			setupMocks: func(repo *MockAccountRepository, cache *MockCacheService, logger *MockLogger) {
				repo.On("GetByAccountName", mock.Anything, "Test Account").Return(nil, errs.ErrAccountNotFound)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Account")).Return(nil)
				cache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.Anything, 15*time.Minute).Return(nil)
				logger.On("Info", mock.Anything, mock.Anything).Return()
				logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
			},
			expectedError: nil,
			validateResult: func(t *testing.T, result *dto.AccountResponse) {
				assert.Equal(t, "Test Account", result.AccountName)
				assert.Equal(t, 1000.0, result.Balance)
				assert.Equal(t, "ACTIVE", result.Status)
			},
		},
		{
			name: "fail_account_already_exists",
			request: dto.CreateAccountRequest{
				AccountName:    "Existing Account",
				InitialBalance: 500.0,
			},
			setupMocks: func(repo *MockAccountRepository, cache *MockCacheService, logger *MockLogger) {
				existingAccount := createTestAccount()
				repo.On("GetByAccountName", mock.Anything, "Existing Account").Return(existingAccount, nil)
				logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
				logger.On("Warn", mock.Anything, mock.Anything).Return()
			},
			expectedError: errs.ErrAccountAlreadyExists,
			validateResult: func(t *testing.T, result *dto.AccountResponse) {
				assert.Nil(t, result)
			},
		},
		{
			name: "fail_repository_error",
			request: dto.CreateAccountRequest{
				AccountName:    "Test Account",
				InitialBalance: 1000.0,
			},
			setupMocks: func(repo *MockAccountRepository, cache *MockCacheService, logger *MockLogger) {
				repo.On("GetByAccountName", mock.Anything, "Test Account").Return(nil, errs.ErrAccountNotFound)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Account")).Return(errors.New("database error"))
				logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
				logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
			},
			expectedError: errors.New("database error"),
			validateResult: func(t *testing.T, result *dto.AccountResponse) {
				assert.Nil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockRepo := new(MockAccountRepository)
			mockCache := new(MockCacheService)
			mockLogger := new(MockLogger)

			tt.setupMocks(mockRepo, mockCache, mockLogger)

			// Create use case
			uc := NewAccountUseCase(mockRepo, mockCache, mockLogger)

			// Execute
			result, err := uc.CreateAccount(context.Background(), tt.request)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			tt.validateResult(t, result)

			// Verify mocks
			mockRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}

func TestAccountUseCase_GetAccount(t *testing.T) {
	tests := []struct {
		name           string
		accountID      string
		setupMocks     func(*MockAccountRepository, *MockCacheService, *MockLogger)
		expectedError  error
		validateResult func(*testing.T, *dto.AccountResponse)
	}{
		{
			name:      "success_get_from_repository",
			accountID: "2024072912345678",
			setupMocks: func(repo *MockAccountRepository, cache *MockCacheService, logger *MockLogger) {
				account := createTestAccount()
				cache.On("Get", mock.Anything, "account:2024072912345678", mock.Anything).Return(errors.New("cache miss"))
				repo.On("GetByID", mock.Anything, mock.AnythingOfType("vo.AccountID")).Return(account, nil)
				cache.On("Set", mock.Anything, "account:2024072912345678", mock.Anything, 15*time.Minute).Return(nil)
				logger.On("Debug", mock.Anything, mock.Anything).Return()
				logger.On("Debug", mock.Anything, mock.Anything, mock.Anything).Return()
			},
			expectedError: nil,
			validateResult: func(t *testing.T, result *dto.AccountResponse) {
				assert.NotNil(t, result)
				assert.Equal(t, "Test Account", result.AccountName)
			},
		},
		{
			name:      "fail_invalid_account_id",
			accountID: "invalid-id",
			setupMocks: func(repo *MockAccountRepository, cache *MockCacheService, logger *MockLogger) {
				logger.On("Debug", mock.Anything, mock.Anything).Return()
				logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
			},
			expectedError: errs.ErrInvalidAccountID,
			validateResult: func(t *testing.T, result *dto.AccountResponse) {
				assert.Nil(t, result)
			},
		},
		{
			name:      "fail_account_not_found",
			accountID: "2024072912345678",
			setupMocks: func(repo *MockAccountRepository, cache *MockCacheService, logger *MockLogger) {
				cache.On("Get", mock.Anything, "account:2024072912345678", mock.Anything).Return(errors.New("cache miss"))
				repo.On("GetByID", mock.Anything, mock.AnythingOfType("vo.AccountID")).Return(&entity.Account{}, errs.ErrAccountNotFound)
				logger.On("Debug", mock.Anything, mock.Anything).Return()
				logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
			},
			expectedError: errs.ErrAccountNotFound,
			validateResult: func(t *testing.T, result *dto.AccountResponse) {
				assert.Nil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockRepo := new(MockAccountRepository)
			mockCache := new(MockCacheService)
			mockLogger := new(MockLogger)

			tt.setupMocks(mockRepo, mockCache, mockLogger)

			// Create use case
			uc := NewAccountUseCase(mockRepo, mockCache, mockLogger)

			// Execute
			result, err := uc.GetAccount(context.Background(), tt.accountID)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			tt.validateResult(t, result)

			// Verify mocks
			mockRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}

func TestAccountUseCase_UpdateAccount(t *testing.T) {
	tests := []struct {
		name           string
		request        dto.UpdateAccountRequest
		setupMocks     func(*MockAccountRepository, *MockCacheService, *MockLogger)
		expectedError  error
		validateResult func(*testing.T, *dto.AccountResponse)
	}{
		{
			name: "success_update_account",
			request: dto.UpdateAccountRequest{
				ID:          "2024072912345678",
				AccountName: "Updated Account Name",
			},
			setupMocks: func(repo *MockAccountRepository, cache *MockCacheService, logger *MockLogger) {
				account := createTestAccount()
				repo.On("GetByID", mock.Anything, mock.AnythingOfType("vo.AccountID")).Return(account, nil)
				repo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Account")).Return(nil)
				cache.On("Set", mock.Anything, "account:2024072912345678", mock.Anything, 15*time.Minute).Return(nil)
				logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
			},
			expectedError: nil,
			validateResult: func(t *testing.T, result *dto.AccountResponse) {
				assert.NotNil(t, result)
				assert.Equal(t, "Updated Account Name", result.AccountName)
			},
		},
		{
			name: "fail_account_not_found",
			request: dto.UpdateAccountRequest{
				ID:          "2024072912345678",
				AccountName: "Updated Account Name",
			},
			setupMocks: func(repo *MockAccountRepository, cache *MockCacheService, logger *MockLogger) {
				repo.On("GetByID", mock.Anything, mock.AnythingOfType("vo.AccountID")).Return(&entity.Account{}, errs.ErrAccountNotFound)
				logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
				logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
			},
			expectedError: errs.ErrAccountNotFound,
			validateResult: func(t *testing.T, result *dto.AccountResponse) {
				assert.Nil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockRepo := new(MockAccountRepository)
			mockCache := new(MockCacheService)
			mockLogger := new(MockLogger)

			tt.setupMocks(mockRepo, mockCache, mockLogger)

			// Create use case
			uc := NewAccountUseCase(mockRepo, mockCache, mockLogger)

			// Execute
			result, err := uc.UpdateAccount(context.Background(), tt.request)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			tt.validateResult(t, result)

			// Verify mocks
			mockRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}

func TestAccountUseCase_DeleteAccount(t *testing.T) {
	tests := []struct {
		name          string
		accountID     string
		setupMocks    func(*MockAccountRepository, *MockCacheService, *MockLogger)
		expectedError error
	}{
		{
			name:      "success_delete_account",
			accountID: "2024072912345678",
			setupMocks: func(repo *MockAccountRepository, cache *MockCacheService, logger *MockLogger) {
				account := createTestAccount()
				repo.On("GetByID", mock.Anything, mock.AnythingOfType("vo.AccountID")).Return(account, nil)
				repo.On("Delete", mock.Anything, mock.AnythingOfType("vo.AccountID")).Return(nil)
				cache.On("Delete", mock.Anything, "account:2024072912345678").Return(nil)
				logger.On("Info", mock.Anything, mock.Anything).Return()
				logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
			},
			expectedError: nil,
		},
		{
			name:      "fail_account_not_found",
			accountID: "2024072912345678",
			setupMocks: func(repo *MockAccountRepository, cache *MockCacheService, logger *MockLogger) {
				repo.On("GetByID", mock.Anything, mock.AnythingOfType("vo.AccountID")).Return(&entity.Account{}, errs.ErrAccountNotFound)
				logger.On("Info", mock.Anything, mock.Anything).Return()
				logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
			},
			expectedError: errs.ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockRepo := new(MockAccountRepository)
			mockCache := new(MockCacheService)
			mockLogger := new(MockLogger)

			tt.setupMocks(mockRepo, mockCache, mockLogger)

			// Create use case
			uc := NewAccountUseCase(mockRepo, mockCache, mockLogger)

			// Execute
			err := uc.DeleteAccount(context.Background(), tt.accountID)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify mocks
			mockRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}

func TestAccountUseCase_SuspendAccount(t *testing.T) {
	tests := []struct {
		name          string
		accountID     string
		setupMocks    func(*MockAccountRepository, *MockCacheService, *MockLogger)
		expectedError error
	}{
		{
			name:      "success_suspend_account",
			accountID: "2024072912345678",
			setupMocks: func(repo *MockAccountRepository, cache *MockCacheService, logger *MockLogger) {
				account := createTestAccount()
				repo.On("GetByID", mock.Anything, mock.AnythingOfType("vo.AccountID")).Return(account, nil)
				repo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Account")).Return(nil)
				cache.On("Set", mock.Anything, "account:2024072912345678", mock.Anything, 15*time.Minute).Return(nil)
				logger.On("Info", mock.Anything, mock.Anything).Return()
				logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
			},
			expectedError: nil,
		},
		{
			name:      "fail_account_not_found",
			accountID: "2024072912345678",
			setupMocks: func(repo *MockAccountRepository, cache *MockCacheService, logger *MockLogger) {
				repo.On("GetByID", mock.Anything, mock.AnythingOfType("vo.AccountID")).Return(&entity.Account{}, errs.ErrAccountNotFound)
				logger.On("Info", mock.Anything, mock.Anything).Return()
				logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
			},
			expectedError: errs.ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockRepo := new(MockAccountRepository)
			mockCache := new(MockCacheService)
			mockLogger := new(MockLogger)

			tt.setupMocks(mockRepo, mockCache, mockLogger)

			// Create use case
			uc := NewAccountUseCase(mockRepo, mockCache, mockLogger)

			// Execute
			err := uc.SuspendAccount(context.Background(), tt.accountID)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify mocks
			mockRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}

func TestAccountUseCase_ActivateAccount(t *testing.T) {
	tests := []struct {
		name          string
		accountID     string
		setupMocks    func(*MockAccountRepository, *MockCacheService, *MockLogger)
		expectedError error
	}{
		{
			name:      "success_activate_account",
			accountID: "2024072912345678",
			setupMocks: func(repo *MockAccountRepository, cache *MockCacheService, logger *MockLogger) {
				account := createTestAccount()
				account.Status = vo.AccountStatusSuspended // Set to suspended so it can be activated
				repo.On("GetByID", mock.Anything, mock.AnythingOfType("vo.AccountID")).Return(account, nil)
				repo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Account")).Return(nil)
				cache.On("Set", mock.Anything, "account:2024072912345678", mock.Anything, 15*time.Minute).Return(nil)
				logger.On("Info", mock.Anything, mock.Anything).Return()
				logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
			},
			expectedError: nil,
		},
		{
			name:      "fail_account_not_found",
			accountID: "2024072912345678",
			setupMocks: func(repo *MockAccountRepository, cache *MockCacheService, logger *MockLogger) {
				repo.On("GetByID", mock.Anything, mock.AnythingOfType("vo.AccountID")).Return(&entity.Account{}, errs.ErrAccountNotFound)
				logger.On("Info", mock.Anything, mock.Anything).Return()
				logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
			},
			expectedError: errs.ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockRepo := new(MockAccountRepository)
			mockCache := new(MockCacheService)
			mockLogger := new(MockLogger)

			tt.setupMocks(mockRepo, mockCache, mockLogger)

			// Create use case
			uc := NewAccountUseCase(mockRepo, mockCache, mockLogger)

			// Execute
			err := uc.ActivateAccount(context.Background(), tt.accountID)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify mocks
			mockRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}
