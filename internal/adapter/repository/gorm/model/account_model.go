package model

import (
	"time"

	"github.com/hydr0g3nz/mini_bank/internal/domain/entity"
	"github.com/hydr0g3nz/mini_bank/internal/domain/vo"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	AccountID   string          `gorm:"size:16;uniqueIndex;not null"` // Format: YYYYMMDD + 8 digits
	AccountName string          `gorm:"size:100;not null"`
	Balance     decimal.Decimal `gorm:"type:decimal(20,2);not null;default:0"`
	Status      string          `gorm:"size:20;not null;default:'ACTIVE'"` // ACTIVE, INACTIVE, SUSPENDED
	CreatedAt   time.Time       `gorm:"not null"`
	UpdatedAt   time.Time       `gorm:"not null"`
}

// TableName specifies the table name for the Account model
func (Account) TableName() string {
	return "accounts"
}

// ToDomainAccount converts GORM model to domain entity
func (a *Account) ToDomainAccount() (*entity.Account, error) {
	accountID, err := vo.NewAccountIDFromString(a.AccountID)
	if err != nil {
		return nil, err
	}

	money := vo.NewMoney(a.Balance)
	status := vo.AccountStatus(a.Status)

	return &entity.Account{
		ID:          accountID,
		AccountName: a.AccountName,
		Balance:     money,
		Status:      status,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}, nil
}

// FromDomainAccount converts domain entity to GORM model
func FromDomainAccount(domainAccount *entity.Account) *Account {
	return &Account{
		Model: gorm.Model{
			ID:        uint(0), // Will be auto-generated
			CreatedAt: domainAccount.CreatedAt,
			UpdatedAt: domainAccount.UpdatedAt,
		},
		AccountID:   domainAccount.ID.String(),
		AccountName: domainAccount.AccountName,
		Balance:     domainAccount.Balance.Amount(),
		Status:      string(domainAccount.Status),
	}
}

// UpdateFromDomain updates the GORM model with domain entity data (preserves GORM ID)
func (a *Account) UpdateFromDomain(domainAccount *entity.Account) {
	a.AccountID = domainAccount.ID.String()
	a.AccountName = domainAccount.AccountName
	a.Balance = domainAccount.Balance.Amount()
	a.Status = string(domainAccount.Status)
	a.UpdatedAt = domainAccount.UpdatedAt
}
