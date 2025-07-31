package model

import (
	"time"

	"github.com/hydr0g3nz/mini_bank/internal/domain/entity"
	"github.com/hydr0g3nz/mini_bank/internal/domain/vo"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	TransactionID   string          `gorm:"size:25;uniqueIndex;not null"` // Format: TXN + timestamp + random
	FromAccountID   *string         `gorm:"size:16;index"`                // Foreign key to accounts.account_id
	ToAccountID     *string         `gorm:"size:16;index"`                // Foreign key to accounts.account_id
	TransactionType string          `gorm:"size:20;not null"`             // DEBIT, CREDIT, TRANSFER
	Amount          decimal.Decimal `gorm:"type:decimal(20,2);not null"`
	Description     string          `gorm:"size:500"`
	Reference       string          `gorm:"size:100"`
	Status          string          `gorm:"size:20;not null;default:'PENDING'"` // PENDING, COMPLETED, FAILED, CANCELLED
	CreatedAt       time.Time       `gorm:"not null"`
	CompletedAt     *time.Time      `gorm:"index"`
}

// TableName specifies the table name for the Transaction model
func (Transaction) TableName() string {
	return "transactions"
}

// ToDomainTransaction converts GORM model to domain entity
func (t *Transaction) ToDomainTransaction() (*entity.Transaction, error) {
	transactionID, err := vo.NewTransactionIDFromString(t.TransactionID)
	if err != nil {
		return nil, err
	}

	var fromAccountID *vo.AccountID
	if t.FromAccountID != nil {
		fromID, err := vo.NewAccountIDFromString(*t.FromAccountID)
		if err != nil {
			return nil, err
		}
		fromAccountID = &fromID
	}

	var toAccountID *vo.AccountID
	if t.ToAccountID != nil {
		toID, err := vo.NewAccountIDFromString(*t.ToAccountID)
		if err != nil {
			return nil, err
		}
		toAccountID = &toID
	}

	money := vo.NewMoney(t.Amount)
	transactionType := vo.TransactionType(t.TransactionType)
	status := vo.TransactionStatus(t.Status)

	return &entity.Transaction{
		ID:              transactionID,
		FromAccountID:   fromAccountID,
		ToAccountID:     toAccountID,
		TransactionType: transactionType,
		Amount:          money,
		Description:     t.Description,
		Reference:       t.Reference,
		Status:          status,
		CreatedAt:       t.CreatedAt,
		CompletedAt:     t.CompletedAt,
	}, nil
}

// FromDomainTransaction converts domain entity to GORM model
func FromDomainTransaction(domainTransaction *entity.Transaction) *Transaction {
	var fromAccountID *string
	if domainTransaction.FromAccountID != nil {
		id := domainTransaction.FromAccountID.String()
		fromAccountID = &id
	}

	var toAccountID *string
	if domainTransaction.ToAccountID != nil {
		id := domainTransaction.ToAccountID.String()
		toAccountID = &id
	}

	return &Transaction{
		Model: gorm.Model{
			ID:        uint(0), // Will be auto-generated
			CreatedAt: domainTransaction.CreatedAt,
		},
		TransactionID:   domainTransaction.ID.String(),
		FromAccountID:   fromAccountID,
		ToAccountID:     toAccountID,
		TransactionType: string(domainTransaction.TransactionType),
		Amount:          domainTransaction.Amount.Amount(),
		Description:     domainTransaction.Description,
		Reference:       domainTransaction.Reference,
		Status:          string(domainTransaction.Status),
		CompletedAt:     domainTransaction.CompletedAt,
	}
}

// UpdateFromDomain updates the GORM model with domain entity data (preserves GORM ID)
func (t *Transaction) UpdateFromDomain(domainTransaction *entity.Transaction) {
	t.TransactionID = domainTransaction.ID.String()

	var fromAccountID *string
	if domainTransaction.FromAccountID != nil {
		id := domainTransaction.FromAccountID.String()
		fromAccountID = &id
	}
	t.FromAccountID = fromAccountID

	var toAccountID *string
	if domainTransaction.ToAccountID != nil {
		id := domainTransaction.ToAccountID.String()
		toAccountID = &id
	}
	t.ToAccountID = toAccountID

	t.TransactionType = string(domainTransaction.TransactionType)
	t.Amount = domainTransaction.Amount.Amount()
	t.Description = domainTransaction.Description
	t.Reference = domainTransaction.Reference
	t.Status = string(domainTransaction.Status)
	t.CompletedAt = domainTransaction.CompletedAt
	t.UpdatedAt = time.Now()
}
