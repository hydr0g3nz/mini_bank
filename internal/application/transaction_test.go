package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hydr0g3nz/mini_bank/internal/application/dto"
	"github.com/hydr0g3nz/mini_bank/internal/domain/entity"
	errs "github.com/hydr0g3nz/mini_bank/internal/domain/error"
	"github.com/hydr0g3nz/mini_bank/internal/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock Transaction Repository
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetByID(ctx context.Context, id vo.TransactionID) (*entity.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Update(ctx context.Context, transaction *entity.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) List(ctx context.Context, limit, offset int) ([]*entity.Transaction, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByAccountID(ctx context.Context, accountID vo.AccountID, limit, offset int) ([]*entity.Transaction, error) {
	args := m.Called(ctx, accountID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByStatus(ctx context.Context, status vo.TransactionStatus, limit, offset int) ([]*entity.Transaction, error) {
	args := m.Called(ctx, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Transaction), args.Error(1)
}

// Test Suite
type TransactionUseCaseTestSuite struct {
	suite.Suite
	usecase         TransactionUseCase
	mockTxnRepo     *MockTransactionRepository
	mockAccountRepo *MockAccountRepository
	mockCache       *MockCacheService
	mockLogger      *MockLogger
	ctx             context.Context
	testAccount     *entity.Account
	testTransaction *entity.Transaction
}

func (suite *TransactionUseCaseTestSuite) SetupTest() {
	suite.mockTxnRepo = new(MockTransactionRepository)
	suite.mockAccountRepo = new(MockAccountRepository)
	suite.mockCache = new(MockCacheService)
	suite.mockLogger = new(MockLogger)
	suite.ctx = context.Background()

	// Allow logger calls without strict expectations
	suite.mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	suite.mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	suite.mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()
	suite.mockLogger.On("Warn", mock.Anything, mock.Anything).Maybe()

	suite.usecase = NewTransactionUseCase(suite.mockTxnRepo, suite.mockAccountRepo, suite.mockCache, suite.mockLogger).(*transactionUseCase)

	// Create test account
	var err error
	suite.testAccount, err = entity.NewAccount("Test Account", vo.NewMoneyFromFloat(1000.0))
	suite.Require().NoError(err)

	// Create test transaction
	suite.testTransaction, err = entity.NewDebitTransaction(
		suite.testAccount.ID,
		vo.NewMoneyFromFloat(100.0),
		"Test debit",
		"TEST-REF",
	)
	suite.Require().NoError(err)
}

func (suite *TransactionUseCaseTestSuite) TestCreateTransaction_Debit_Success() {
	fromAccountID := suite.testAccount.ID.String()
	req := dto.CreateTransactionRequest{
		FromAccountID:   &fromAccountID,
		TransactionType: "DEBIT",
		Amount:          100.0,
		Description:     "Test debit",
		Reference:       "TEST-REF",
	}

	suite.mockAccountRepo.On("GetByID", suite.ctx, suite.testAccount.ID).Return(suite.testAccount, nil)
	suite.mockTxnRepo.On("Create", suite.ctx, mock.AnythingOfType("*entity.Transaction")).Return(nil)
	suite.mockCache.On("Set", suite.ctx, mock.AnythingOfType("string"), mock.Anything, 30*time.Minute).Return(nil)

	result, err := suite.usecase.CreateTransaction(suite.ctx, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "DEBIT", result.TransactionType)
	assert.Equal(suite.T(), 100.0, result.Amount)
	suite.mockTxnRepo.AssertExpectations(suite.T())
	suite.mockAccountRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestCreateTransaction_Credit_Success() {
	toAccountID := suite.testAccount.ID.String()
	req := dto.CreateTransactionRequest{
		ToAccountID:     &toAccountID,
		TransactionType: "CREDIT",
		Amount:          100.0,
		Description:     "Test credit",
		Reference:       "TEST-REF",
	}

	suite.mockAccountRepo.On("GetByID", suite.ctx, suite.testAccount.ID).Return(suite.testAccount, nil)
	suite.mockTxnRepo.On("Create", suite.ctx, mock.AnythingOfType("*entity.Transaction")).Return(nil)
	suite.mockCache.On("Set", suite.ctx, mock.AnythingOfType("string"), mock.Anything, 30*time.Minute).Return(nil)

	result, err := suite.usecase.CreateTransaction(suite.ctx, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "CREDIT", result.TransactionType)
	assert.Equal(suite.T(), 100.0, result.Amount)
	suite.mockTxnRepo.AssertExpectations(suite.T())
	suite.mockAccountRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestCreateTransaction_Transfer_Success() {
	// Create second account
	toAccount, _ := entity.NewAccount("To Account", vo.NewMoneyFromFloat(500.0))

	fromAccountID := suite.testAccount.ID.String()
	toAccountID := toAccount.ID.String()
	req := dto.CreateTransactionRequest{
		FromAccountID:   &fromAccountID,
		ToAccountID:     &toAccountID,
		TransactionType: "TRANSFER",
		Amount:          100.0,
		Description:     "Test transfer",
		Reference:       "TEST-REF",
	}

	suite.mockAccountRepo.On("GetByID", suite.ctx, suite.testAccount.ID).Return(suite.testAccount, nil)
	suite.mockAccountRepo.On("GetByID", suite.ctx, toAccount.ID).Return(toAccount, nil)
	suite.mockTxnRepo.On("Create", suite.ctx, mock.AnythingOfType("*entity.Transaction")).Return(nil)
	suite.mockCache.On("Set", suite.ctx, mock.AnythingOfType("string"), mock.Anything, 30*time.Minute).Return(nil)

	result, err := suite.usecase.CreateTransaction(suite.ctx, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "TRANSFER", result.TransactionType)
	assert.Equal(suite.T(), 100.0, result.Amount)
	suite.mockTxnRepo.AssertExpectations(suite.T())
	suite.mockAccountRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestCreateTransaction_AccountNotFound() {
	fromAccountID := suite.testAccount.ID.String()
	req := dto.CreateTransactionRequest{
		FromAccountID:   &fromAccountID,
		TransactionType: "DEBIT",
		Amount:          100.0,
		Description:     "Test debit",
		Reference:       "TEST-REF",
	}

	suite.mockAccountRepo.On("GetByID", suite.ctx, suite.testAccount.ID).Return((*entity.Account)(nil), errs.ErrAccountNotFound)

	result, err := suite.usecase.CreateTransaction(suite.ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), errs.ErrAccountNotFound, err)
	suite.mockAccountRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestConfirmTransaction_Success() {
	req := dto.ConfirmTransactionRequest{
		ID: suite.testTransaction.ID.String(),
	}

	// Mock cache miss for idempotency check
	idempotencyKey := "confirm_transaction:" + req.ID
	suite.mockCache.On("Get", suite.ctx, idempotencyKey, mock.Anything).Return(errors.New("cache miss"))

	// Mock lock acquisition
	lockKey := "lock:transaction:" + req.ID
	suite.mockCache.On("Set", suite.ctx, lockKey, mock.Anything, 30*time.Second).Return(nil)
	suite.mockCache.On("Delete", suite.ctx, lockKey).Return(nil)

	// Mock transaction retrieval
	suite.mockTxnRepo.On("GetByID", suite.ctx, suite.testTransaction.ID).Return(suite.testTransaction, nil)

	// Mock account operations for debit transaction
	suite.mockAccountRepo.On("GetByID", suite.ctx, *suite.testTransaction.FromAccountID).Return(suite.testAccount, nil)
	suite.mockAccountRepo.On("Update", suite.ctx, mock.AnythingOfType("*entity.Account")).Return(nil)

	// Mock transaction update
	suite.mockTxnRepo.On("Update", suite.ctx, mock.AnythingOfType("*entity.Transaction")).Return(nil)

	// Mock cache operations
	suite.mockCache.On("Set", suite.ctx, idempotencyKey, mock.Anything, 24*time.Hour).Return(nil)
	suite.mockCache.On("Set", suite.ctx, "transaction:"+req.ID, mock.Anything, 30*time.Minute).Return(nil)
	suite.mockCache.On("Delete", suite.ctx, "account:"+suite.testAccount.ID.String()).Return(nil)

	result, err := suite.usecase.ConfirmTransaction(suite.ctx, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "COMPLETED", result.Status)
	suite.mockTxnRepo.AssertExpectations(suite.T())
	suite.mockAccountRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestConfirmTransaction_AlreadyCompleted() {
	// Create completed transaction
	completedTxn, _ := entity.NewDebitTransaction(
		suite.testAccount.ID,
		vo.NewMoneyFromFloat(100.0),
		"Test debit",
		"TEST-REF",
	)
	completedTxn.MarkAsCompleted()

	req := dto.ConfirmTransactionRequest{
		ID: completedTxn.ID.String(),
	}

	// Mock cache miss for idempotency check
	idempotencyKey := "confirm_transaction:" + req.ID
	suite.mockCache.On("Get", suite.ctx, idempotencyKey, mock.Anything).Return(errors.New("cache miss"))

	// Mock lock acquisition
	lockKey := "lock:transaction:" + req.ID
	suite.mockCache.On("Set", suite.ctx, lockKey, mock.Anything, 30*time.Second).Return(nil)
	suite.mockCache.On("Delete", suite.ctx, lockKey).Return(nil)

	// Mock transaction retrieval
	suite.mockTxnRepo.On("GetByID", suite.ctx, completedTxn.ID).Return(completedTxn, nil)

	// Mock cache set for idempotent result
	suite.mockCache.On("Set", suite.ctx, idempotencyKey, mock.Anything, 24*time.Hour).Return(nil)

	result, err := suite.usecase.ConfirmTransaction(suite.ctx, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "COMPLETED", result.Status)
	suite.mockTxnRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestConfirmTransaction_NotFound() {
	req := dto.ConfirmTransactionRequest{
		ID: suite.testTransaction.ID.String(),
	}

	// Mock cache miss for idempotency check
	idempotencyKey := "confirm_transaction:" + req.ID
	suite.mockCache.On("Get", suite.ctx, idempotencyKey, mock.Anything).Return(errors.New("cache miss"))

	// Mock lock acquisition
	lockKey := "lock:transaction:" + req.ID
	suite.mockCache.On("Set", suite.ctx, lockKey, mock.Anything, 30*time.Second).Return(nil)
	suite.mockCache.On("Delete", suite.ctx, lockKey).Return(nil)

	// Mock transaction not found
	suite.mockTxnRepo.On("GetByID", suite.ctx, suite.testTransaction.ID).Return(nil, errs.ErrTransactionNotFound)

	result, err := suite.usecase.ConfirmTransaction(suite.ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), errs.ErrTransactionNotFound, err)
	suite.mockTxnRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestGetTransaction_Success() {
	transactionID := suite.testTransaction.ID.String()

	suite.mockCache.On("Get", suite.ctx, "transaction:"+transactionID, mock.Anything).Return(errors.New("cache miss"))
	suite.mockTxnRepo.On("GetByID", suite.ctx, suite.testTransaction.ID).Return(suite.testTransaction, nil)
	suite.mockCache.On("Set", suite.ctx, "transaction:"+transactionID, mock.Anything, 30*time.Minute).Return(nil)

	result, err := suite.usecase.GetTransaction(suite.ctx, transactionID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), transactionID, result.ID)
	suite.mockTxnRepo.AssertExpectations(suite.T())
	suite.mockCache.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestGetTransaction_FromCache() {
	transactionID := suite.testTransaction.ID.String()
	cachedResponse := dto.TransactionResponse{
		ID:              transactionID,
		TransactionType: "DEBIT",
		Amount:          100.0,
		Status:          string(vo.TransactionStatusPending),
	}

	suite.mockCache.On("Get", suite.ctx, "transaction:"+transactionID, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(2).(*dto.TransactionResponse)
		*dest = cachedResponse
	})

	result, err := suite.usecase.GetTransaction(suite.ctx, transactionID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), transactionID, result.ID)
	suite.mockCache.AssertExpectations(suite.T())
	// Repo should not be called when cache hit
	suite.mockTxnRepo.AssertNotCalled(suite.T(), "GetByID")
}

func (suite *TransactionUseCaseTestSuite) TestGetTransaction_NotFound() {
	transactionID := suite.testTransaction.ID.String()

	suite.mockCache.On("Get", suite.ctx, "transaction:"+transactionID, mock.Anything).Return(errors.New("cache miss"))
	suite.mockTxnRepo.On("GetByID", suite.ctx, suite.testTransaction.ID).Return(nil, errs.ErrTransactionNotFound)

	result, err := suite.usecase.GetTransaction(suite.ctx, transactionID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), errs.ErrTransactionNotFound, err)
	suite.mockTxnRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestListTransactions_Success() {
	req := dto.ListRequest{
		Page:     1,
		PageSize: 10,
	}

	transactions := []*entity.Transaction{suite.testTransaction}
	cacheKey := "transactions:list:page:1:size:10"

	suite.mockCache.On("Get", suite.ctx, cacheKey, mock.Anything).Return(errors.New("cache miss"))
	suite.mockTxnRepo.On("List", suite.ctx, 10, 0).Return(transactions, nil)
	suite.mockCache.On("Set", suite.ctx, cacheKey, mock.Anything, 2*time.Minute).Return(nil)

	result, err := suite.usecase.ListTransactions(suite.ctx, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result.Transactions, 1)
	assert.Equal(suite.T(), 1, result.Pagination.Page)
	suite.mockTxnRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestGetTransactionsByAccount_Success() {
	accountID := suite.testAccount.ID.String()
	req := dto.ListRequest{
		Page:     1,
		PageSize: 10,
	}

	transactions := []*entity.Transaction{suite.testTransaction}
	cacheKey := "transactions:account:" + accountID + ":page:1:size:10"

	suite.mockCache.On("Get", suite.ctx, cacheKey, mock.Anything).Return(errors.New("cache miss"))
	suite.mockTxnRepo.On("GetByAccountID", suite.ctx, suite.testAccount.ID, 10, 0).Return(transactions, nil)
	suite.mockCache.On("Set", suite.ctx, cacheKey, mock.Anything, 5*time.Minute).Return(nil)

	result, err := suite.usecase.GetTransactionsByAccount(suite.ctx, accountID, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result.Transactions, 1)
	suite.mockTxnRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestCancelTransaction_Success() {
	req := dto.CancelTransactionRequest{
		ID: suite.testTransaction.ID.String(),
	}

	suite.mockTxnRepo.On("GetByID", suite.ctx, suite.testTransaction.ID).Return(suite.testTransaction, nil)
	suite.mockTxnRepo.On("Update", suite.ctx, mock.AnythingOfType("*entity.Transaction")).Return(nil)
	suite.mockCache.On("Set", suite.ctx, "transaction:"+req.ID, mock.Anything, 30*time.Minute).Return(nil)

	err := suite.usecase.CancelTransaction(suite.ctx, req)

	assert.NoError(suite.T(), err)
	suite.mockTxnRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestCancelTransaction_NotFound() {
	req := dto.CancelTransactionRequest{
		ID: suite.testTransaction.ID.String(),
	}

	suite.mockTxnRepo.On("GetByID", suite.ctx, suite.testTransaction.ID).Return(nil, errs.ErrTransactionNotFound)

	err := suite.usecase.CancelTransaction(suite.ctx, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), errs.ErrTransactionNotFound, err)
	suite.mockTxnRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestCancelTransaction_AlreadyCompleted() {
	// Create completed transaction
	completedTxn, _ := entity.NewDebitTransaction(
		suite.testAccount.ID,
		vo.NewMoneyFromFloat(100.0),
		"Test debit",
		"TEST-REF",
	)
	completedTxn.MarkAsCompleted()

	req := dto.CancelTransactionRequest{
		ID: completedTxn.ID.String(),
	}

	suite.mockTxnRepo.On("GetByID", suite.ctx, completedTxn.ID).Return(completedTxn, nil)

	err := suite.usecase.CancelTransaction(suite.ctx, req)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), errs.ErrTransactionCannotBeCancelled.Error())
	suite.mockTxnRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestGetTransactionsByStatus_Success() {
	status := "PENDING"
	req := dto.ListRequest{
		Page:     1,
		PageSize: 10,
	}

	transactions := []*entity.Transaction{suite.testTransaction}
	cacheKey := "transactions:status:PENDING:page:1:size:10"

	suite.mockCache.On("Get", suite.ctx, cacheKey, mock.Anything).Return(errors.New("cache miss"))
	suite.mockTxnRepo.On("GetByStatus", suite.ctx, vo.TransactionStatusPending, 10, 0).Return(transactions, nil)
	suite.mockCache.On("Set", suite.ctx, cacheKey, mock.Anything, 5*time.Minute).Return(nil)

	result, err := suite.usecase.GetTransactionsByStatus(suite.ctx, status, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result.Transactions, 1)
	suite.mockTxnRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestGetTransactionsByStatus_InvalidStatus() {
	status := "INVALID_STATUS"
	req := dto.ListRequest{
		Page:     1,
		PageSize: 10,
	}

	result, err := suite.usecase.GetTransactionsByStatus(suite.ctx, status, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "invalid transaction status")
}

func (suite *TransactionUseCaseTestSuite) TestConfirmTransaction_InsufficientBalance() {
	// Create account with low balance
	lowBalanceAccount, _ := entity.NewAccount("Low Balance Account", vo.NewMoneyFromFloat(50.0))

	// Create transaction with amount higher than balance
	highAmountTxn, _ := entity.NewDebitTransaction(
		lowBalanceAccount.ID,
		vo.NewMoneyFromFloat(100.0),
		"Test debit",
		"TEST-REF",
	)

	req := dto.ConfirmTransactionRequest{
		ID: highAmountTxn.ID.String(),
	}

	// Mock cache miss for idempotency check
	idempotencyKey := "confirm_transaction:" + req.ID
	suite.mockCache.On("Get", suite.ctx, idempotencyKey, mock.Anything).Return(errors.New("cache miss"))

	// Mock lock acquisition
	lockKey := "lock:transaction:" + req.ID
	suite.mockCache.On("Set", suite.ctx, lockKey, mock.Anything, 30*time.Second).Return(nil)
	suite.mockCache.On("Delete", suite.ctx, lockKey).Return(nil)

	// Mock transaction retrieval
	suite.mockTxnRepo.On("GetByID", suite.ctx, highAmountTxn.ID).Return(highAmountTxn, nil)

	// Mock account retrieval with low balance
	suite.mockAccountRepo.On("GetByID", suite.ctx, *highAmountTxn.FromAccountID).Return(lowBalanceAccount, nil)

	// Mock transaction update to failed status
	suite.mockTxnRepo.On("Update", suite.ctx, mock.AnythingOfType("*entity.Transaction")).Return(nil)

	result, err := suite.usecase.ConfirmTransaction(suite.ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), errs.ErrInsufficientBalance, err)
	suite.mockTxnRepo.AssertExpectations(suite.T())
	suite.mockAccountRepo.AssertExpectations(suite.T())
}

func TestTransactionUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionUseCaseTestSuite))
}
