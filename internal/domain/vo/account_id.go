package vo

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"time"

	errs "github.com/hydr0g3nz/mini_bank/internal/domain/error"
)

// AccountID represents an account identifier
// Uses numeric format: YYYYMMDD + 8-digit sequence (e.g., 2024072912345678)
type AccountID struct {
	value string
}

// NewAccountID creates a new AccountID with date prefix + random sequence
func NewAccountID() AccountID {
	now := time.Now()
	datePrefix := now.Format("20060102") // YYYYMMDD format

	// Generate 8-digit random sequence
	max := big.NewInt(99999999)
	n, _ := rand.Int(rand.Reader, max)
	sequence := fmt.Sprintf("%08d", n.Int64())

	return AccountID{value: datePrefix + sequence}
}

// NewAccountIDFromString creates AccountID from string with validation
func NewAccountIDFromString(id string) (AccountID, error) {
	if err := validateAccountID(id); err != nil {
		return AccountID{}, err
	}
	return AccountID{value: id}, nil
}

// String returns string representation
func (id AccountID) String() string {
	return id.value
}

// IsEmpty checks if ID is empty
func (id AccountID) IsEmpty() bool {
	return id.value == ""
}

// IsValid checks if ID format is valid
func (id AccountID) IsValid() bool {
	return validateAccountID(id.value) == nil
}

func validateAccountID(id string) error {
	if id == "" {
		return errs.ErrInvalidAccountID
	}

	// Check length (YYYYMMDD + 8 digits = 16 chars)
	if len(id) != 16 {
		return errs.ErrInvalidAccountID
	}

	// Check if all characters are digits
	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		return errs.ErrInvalidAccountID
	}
	if len(id) < 8 {
		return errs.ErrInvalidAccountID
	}
	// Validate date part (first 8 chars)
	dateStr := id[:8]
	if _, err := time.Parse("20060102", dateStr); err != nil {
		return errs.ErrInvalidAccountID
	}

	return nil
}
