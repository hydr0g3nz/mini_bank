package vo

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	errs "github.com/hydr0g3nz/mini_bank/internal/domain/error"
)

// TransactionID represents a transaction identifier
// Format: TXN + timestamp + random suffix (e.g., TXN20240729143045001234)
type TransactionID struct {
	value string
}

// NewTransactionID creates a new TransactionID
func NewTransactionID() TransactionID {
	now := time.Now()
	timestamp := now.Format("20060102150405") // YYYYMMDDHHmmss

	// Generate 6-digit random suffix
	max := big.NewInt(999999)
	n, _ := rand.Int(rand.Reader, max)
	suffix := fmt.Sprintf("%06d", n.Int64())

	return TransactionID{value: "TXN" + timestamp + suffix}
}

// NewTransactionIDFromString creates TransactionID from string with validation
func NewTransactionIDFromString(id string) (TransactionID, error) {
	if err := validateTransactionID(id); err != nil {
		return TransactionID{}, err
	}
	return TransactionID{value: id}, nil
}

// MustNewTransactionIDFromString creates TransactionID from string, panics on error
func MustNewTransactionIDFromString(id string) TransactionID {
	txnID, err := NewTransactionIDFromString(id)
	if err != nil {
		panic(err)
	}
	return txnID
}

// String returns string representation
func (id TransactionID) String() string {
	return id.value
}

// IsEmpty checks if ID is empty
func (id TransactionID) IsEmpty() bool {
	return id.value == ""
}

// IsValid checks if ID format is valid
func (id TransactionID) IsValid() bool {
	return validateTransactionID(id.value) == nil
}

func validateTransactionID(id string) error {
	if id == "" {
		return errs.ErrInvalidTransactionID
	}

	// Must start with "TXN"
	if !strings.HasPrefix(id, "TXN") {
		return errs.ErrInvalidTransactionID
	}

	// Check minimum length (TXN + 14 chars timestamp + 6 chars suffix = 23)
	if len(id) < 23 {
		return errs.ErrInvalidTransactionID
	}

	// Validate timestamp part (chars 3-16)
	if len(id) >= 17 {
		timestampStr := id[3:17]
		if _, err := time.Parse("20060102150405", timestampStr); err != nil {
			return errs.ErrInvalidTransactionID
		}
	}

	// Check if suffix is numeric (chars 17 onwards)
	if len(id) > 17 {
		suffix := id[17:]
		if _, err := strconv.ParseInt(suffix, 10, 64); err != nil {
			return errs.ErrInvalidTransactionID
		}
	}

	return nil
}
