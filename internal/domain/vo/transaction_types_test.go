package vo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		txnType  TransactionType
		expected bool
	}{
		{
			name:     "Debit type is valid",
			txnType:  TransactionTypeDebit,
			expected: true,
		},
		{
			name:     "Credit type is valid",
			txnType:  TransactionTypeCredit,
			expected: true,
		},
		{
			name:     "Transfer type is valid",
			txnType:  TransactionTypeTransfer,
			expected: true,
		},
		{
			name:     "Invalid type",
			txnType:  TransactionType("INVALID"),
			expected: false,
		},
		{
			name:     "Empty type",
			txnType:  TransactionType(""),
			expected: false,
		},
		{
			name:     "Random string type",
			txnType:  TransactionType("RANDOM"),
			expected: false,
		},
		{
			name:     "Lowercase type",
			txnType:  TransactionType("debit"),
			expected: false,
		},
		{
			name:     "Mixed case type",
			txnType:  TransactionType("Debit"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.txnType.IsValid())
		})
	}
}

func TestTransactionType_IsDebit(t *testing.T) {
	tests := []struct {
		name     string
		txnType  TransactionType
		expected bool
	}{
		{
			name:     "Debit type",
			txnType:  TransactionTypeDebit,
			expected: true,
		},
		{
			name:     "Credit type",
			txnType:  TransactionTypeCredit,
			expected: false,
		},
		{
			name:     "Transfer type",
			txnType:  TransactionTypeTransfer,
			expected: false,
		},
		{
			name:     "Invalid type",
			txnType:  TransactionType("INVALID"),
			expected: false,
		},
		{
			name:     "Empty type",
			txnType:  TransactionType(""),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.txnType.IsDebit())
		})
	}
}

func TestTransactionType_IsCredit(t *testing.T) {
	tests := []struct {
		name     string
		txnType  TransactionType
		expected bool
	}{
		{
			name:     "Credit type",
			txnType:  TransactionTypeCredit,
			expected: true,
		},
		{
			name:     "Debit type",
			txnType:  TransactionTypeDebit,
			expected: false,
		},
		{
			name:     "Transfer type",
			txnType:  TransactionTypeTransfer,
			expected: false,
		},
		{
			name:     "Invalid type",
			txnType:  TransactionType("INVALID"),
			expected: false,
		},
		{
			name:     "Empty type",
			txnType:  TransactionType(""),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.txnType.IsCredit())
		})
	}
}

func TestTransactionType_IsTransfer(t *testing.T) {
	tests := []struct {
		name     string
		txnType  TransactionType
		expected bool
	}{
		{
			name:     "Transfer type",
			txnType:  TransactionTypeTransfer,
			expected: true,
		},
		{
			name:     "Debit type",
			txnType:  TransactionTypeDebit,
			expected: false,
		},
		{
			name:     "Credit type",
			txnType:  TransactionTypeCredit,
			expected: false,
		},
		{
			name:     "Invalid type",
			txnType:  TransactionType("INVALID"),
			expected: false,
		},
		{
			name:     "Empty type",
			txnType:  TransactionType(""),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.txnType.IsTransfer())
		})
	}
}

func TestTransactionType_Constants(t *testing.T) {
	// Ensure constants have expected string values
	assert.Equal(t, "DEBIT", string(TransactionTypeDebit))
	assert.Equal(t, "CREDIT", string(TransactionTypeCredit))
	assert.Equal(t, "TRANSFER", string(TransactionTypeTransfer))
}

func TestTransactionType_AllValidTypes(t *testing.T) {
	// Test all valid types
	validTypes := []TransactionType{
		TransactionTypeDebit,
		TransactionTypeCredit,
		TransactionTypeTransfer,
	}

	for _, txnType := range validTypes {
		t.Run(string(txnType), func(t *testing.T) {
			assert.True(t, txnType.IsValid())
		})
	}
}

func TestTransactionType_ExclusiveMethods(t *testing.T) {
	// Test that each type method is exclusive
	tests := []struct {
		name       string
		txnType    TransactionType
		isDebit    bool
		isCredit   bool
		isTransfer bool
	}{
		{
			name:       "Debit type exclusivity",
			txnType:    TransactionTypeDebit,
			isDebit:    true,
			isCredit:   false,
			isTransfer: false,
		},
		{
			name:       "Credit type exclusivity",
			txnType:    TransactionTypeCredit,
			isDebit:    false,
			isCredit:   true,
			isTransfer: false,
		},
		{
			name:       "Transfer type exclusivity",
			txnType:    TransactionTypeTransfer,
			isDebit:    false,
			isCredit:   false,
			isTransfer: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isDebit, tt.txnType.IsDebit())
			assert.Equal(t, tt.isCredit, tt.txnType.IsCredit())
			assert.Equal(t, tt.isTransfer, tt.txnType.IsTransfer())
		})
	}
}

func TestTransactionType_InvalidTypes(t *testing.T) {
	// Test various invalid types
	invalidTypes := []TransactionType{
		TransactionType(""),
		TransactionType("INVALID"),
		TransactionType("debit"),      // lowercase
		TransactionType("credit"),     // lowercase
		TransactionType("transfer"),   // lowercase
		TransactionType("Debit"),      // mixed case
		TransactionType("WITHDRAWAL"), // similar but different
		TransactionType("DEPOSIT"),    // similar but different
		TransactionType("PAYMENT"),    // similar but different
		TransactionType("REFUND"),     // similar but different
		TransactionType(" DEBIT "),    // with spaces
		TransactionType("DEBIT "),     // trailing space
		TransactionType(" DEBIT"),     // leading space
	}

	for _, txnType := range invalidTypes {
		t.Run("Invalid_"+string(txnType), func(t *testing.T) {
			assert.False(t, txnType.IsValid())
			assert.False(t, txnType.IsDebit())
			assert.False(t, txnType.IsCredit())
			assert.False(t, txnType.IsTransfer())
		})
	}
}

func TestTransactionType_CaseSensitivity(t *testing.T) {
	// Test that the type checks are case-sensitive
	tests := []struct {
		name    string
		txnType string
		valid   bool
	}{
		{"Uppercase DEBIT", "DEBIT", true},
		{"Lowercase debit", "debit", false},
		{"Mixed case Debit", "Debit", false},
		{"Uppercase CREDIT", "CREDIT", true},
		{"Lowercase credit", "credit", false},
		{"Mixed case Credit", "Credit", false},
		{"Uppercase TRANSFER", "TRANSFER", true},
		{"Lowercase transfer", "transfer", false},
		{"Mixed case Transfer", "Transfer", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txnType := TransactionType(tt.txnType)
			assert.Equal(t, tt.valid, txnType.IsValid())
		})
	}
}

func TestTransactionType_CompleteCoverage(t *testing.T) {
	// Ensure all constants are covered in IsValid method
	allConstants := []TransactionType{
		TransactionTypeDebit,
		TransactionTypeCredit,
		TransactionTypeTransfer,
	}

	for _, constant := range allConstants {
		t.Run("Coverage_"+string(constant), func(t *testing.T) {
			// Each constant should be valid
			assert.True(t, constant.IsValid())

			// And should match exactly one type check
			typeChecks := []bool{
				constant.IsDebit(),
				constant.IsCredit(),
				constant.IsTransfer(),
			}

			// Count true values
			trueCount := 0
			for _, check := range typeChecks {
				if check {
					trueCount++
				}
			}

			// Exactly one should be true
			assert.Equal(t, 1, trueCount, "Each transaction type should match exactly one type check")
		})
	}
}
