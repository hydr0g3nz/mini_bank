package vo

import (
	"strings"
	"testing"
	"time"

	errs "github.com/hydr0g3nz/mini_bank/internal/domain/error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTransactionID(t *testing.T) {
	id := NewTransactionID()

	assert.NotEmpty(t, id.String())
	assert.True(t, strings.HasPrefix(id.String(), "TXN"))
	assert.True(t, len(id.String()) >= 23) // TXN + 14 chars timestamp + 6 chars suffix
	assert.True(t, id.IsValid())
	assert.False(t, id.IsEmpty())

	// Check timestamp format (chars 3-16: YYYYMMDDHHmmss)
	timestampStr := id.String()[3:17]
	_, err := time.Parse("20060102150405", timestampStr)
	assert.NoError(t, err)

	// Check that suffix is numeric
	suffix := id.String()[17:]
	assert.Equal(t, 6, len(suffix))
	for _, char := range suffix {
		assert.True(t, char >= '0' && char <= '9', "Non-numeric character in suffix: %c", char)
	}
}

func TestNewTransactionIDFromString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorType   error
	}{
		{
			name:        "Valid transaction ID",
			input:       "TXN20240729143045123456",
			expectError: false,
		},
		{
			name:        "Empty string",
			input:       "",
			expectError: true,
			errorType:   errs.ErrInvalidTransactionID,
		},
		{
			name:        "Missing TXN prefix",
			input:       "20240729143045123456",
			expectError: true,
			errorType:   errs.ErrInvalidTransactionID,
		},
		{
			name:        "Wrong prefix",
			input:       "ABC20240729143045123456",
			expectError: true,
			errorType:   errs.ErrInvalidTransactionID,
		},
		{
			name:        "Too short",
			input:       "TXN2024072914304512",
			expectError: true,
			errorType:   errs.ErrInvalidTransactionID,
		},
		{
			name:        "Invalid timestamp - month 13",
			input:       "TXN20241329143045123456",
			expectError: true,
			errorType:   errs.ErrInvalidTransactionID,
		},
		{
			name:        "Invalid timestamp - hour 25",
			input:       "TXN20240729253045123456",
			expectError: true,
			errorType:   errs.ErrInvalidTransactionID,
		},
		{
			name:        "Non-numeric suffix",
			input:       "TXN2024072914304512345A",
			expectError: true,
			errorType:   errs.ErrInvalidTransactionID,
		},
		{
			name:        "Valid with longer suffix",
			input:       "TXN202407291430451234567890",
			expectError: false,
		},
		{
			name:        "Current timestamp format",
			input:       "TXN" + time.Now().Format("20060102150405") + "123456",
			expectError: false,
		},
		{
			name:        "Valid edge case timestamps",
			input:       "TXN20241231235959999999",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewTransactionIDFromString(tt.input)

			if tt.expectError {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.errorType)
				assert.True(t, id.IsEmpty())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.input, id.String())
				assert.True(t, id.IsValid())
				assert.False(t, id.IsEmpty())
			}
		})
	}
}

func TestMustNewTransactionIDFromString(t *testing.T) {
	t.Run("Valid ID does not panic", func(t *testing.T) {
		validID := "TXN20240729143045123456"

		assert.NotPanics(t, func() {
			id := MustNewTransactionIDFromString(validID)
			assert.Equal(t, validID, id.String())
		})
	})

	t.Run("Invalid ID panics", func(t *testing.T) {
		invalidID := "invalid"

		assert.Panics(t, func() {
			MustNewTransactionIDFromString(invalidID)
		})
	})
}

func TestTransactionID_String(t *testing.T) {
	validID := "TXN20240729143045123456"
	id, err := NewTransactionIDFromString(validID)
	require.NoError(t, err)

	assert.Equal(t, validID, id.String())
}

func TestTransactionID_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		id       TransactionID
		expected bool
	}{
		{
			name:     "Empty transaction ID",
			id:       TransactionID{},
			expected: true,
		},
		{
			name:     "Valid transaction ID",
			id:       TransactionID{value: "TXN20240729143045123456"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.id.IsEmpty())
		})
	}
}

func TestTransactionID_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		id       TransactionID
		expected bool
	}{
		{
			name:     "Valid transaction ID",
			id:       TransactionID{value: "TXN20240729143045123456"},
			expected: true,
		},
		{
			name:     "Empty transaction ID",
			id:       TransactionID{},
			expected: false,
		},
		{
			name:     "Invalid format - no prefix",
			id:       TransactionID{value: "20240729143045123456"},
			expected: false,
		},
		{
			name:     "Invalid format - wrong prefix",
			id:       TransactionID{value: "ABC20240729143045123456"},
			expected: false,
		},
		{
			name:     "Invalid format - too short",
			id:       TransactionID{value: "TXN123"},
			expected: false,
		},
		{
			name:     "Invalid timestamp",
			id:       TransactionID{value: "TXN20241329143045123456"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.id.IsValid())
		})
	}
}

func TestValidateTransactionID(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "Valid ID",
			input:       "TXN20240729143045123456",
			expectError: false,
		},
		{
			name:        "Empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "Missing TXN prefix",
			input:       "20240729143045123456",
			expectError: true,
		},
		{
			name:        "Wrong prefix",
			input:       "ABC20240729143045123456",
			expectError: true,
		},
		{
			name:        "Too short",
			input:       "TXN123",
			expectError: true,
		},
		{
			name:        "Invalid timestamp - invalid month",
			input:       "TXN20241329143045123456",
			expectError: true,
		},
		{
			name:        "Invalid timestamp - invalid day",
			input:       "TXN20240732143045123456",
			expectError: true,
		},
		{
			name:        "Invalid timestamp - invalid hour",
			input:       "TXN20240729253045123456",
			expectError: true,
		},
		{
			name:        "Invalid timestamp - invalid minute",
			input:       "TXN20240729146045123456",
			expectError: true,
		},
		{
			name:        "Invalid timestamp - invalid second",
			input:       "TXN20240729143065123456",
			expectError: true,
		},
		{
			name:        "Non-numeric suffix",
			input:       "TXN2024072914304512345A",
			expectError: true,
		},
		{
			name:        "Valid edge cases",
			input:       "TXN20241231235959999999",
			expectError: false,
		},
		{
			name:        "Valid leap year timestamp",
			input:       "TXN20240229120000123456",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTransactionID(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.ErrorIs(t, err, errs.ErrInvalidTransactionID)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTransactionID_Uniqueness(t *testing.T) {
	// Generate multiple IDs and ensure they are unique
	ids := make(map[string]bool)
	numIDs := 1000

	for i := 0; i < numIDs; i++ {
		id := NewTransactionID()
		idStr := id.String()

		// Check that this ID hasn't been generated before
		assert.False(t, ids[idStr], "Duplicate ID generated: %s", idStr)
		ids[idStr] = true

		// Verify format
		assert.True(t, strings.HasPrefix(idStr, "TXN"))
		assert.True(t, len(idStr) >= 23)
		assert.True(t, id.IsValid())
	}
}

func TestTransactionID_TimestampPart(t *testing.T) {
	before := time.Now()
	id := NewTransactionID()
	after := time.Now()

	idStr := id.String()
	timestampStr := idStr[3:17] // YYYYMMDDHHmmss

	parsedTime, err := time.ParseInLocation("20060102150405", timestampStr, time.Now().Location())
	require.NoError(t, err)

	acceptableStart := before.Add(-1 * time.Second)
	acceptableEnd := after.Add(1 * time.Second)

	assert.True(t, parsedTime.After(acceptableStart) && parsedTime.Before(acceptableEnd))
}

func TestTransactionID_SuffixPart(t *testing.T) {
	id := NewTransactionID()
	idStr := id.String()

	// Extract suffix part (everything after timestamp)
	suffix := idStr[17:]

	// Should be exactly 6 digits
	assert.Equal(t, 6, len(suffix))

	// Should be all numeric
	for _, char := range suffix {
		assert.True(t, char >= '0' && char <= '9', "Non-numeric character in suffix: %c", char)
	}

	// Test multiple generations to ensure proper zero-padding
	for i := 0; i < 10; i++ {
		testID := NewTransactionID()
		testSuffix := testID.String()[17:]
		assert.Equal(t, 6, len(testSuffix))
		// Verify it's numeric
		for _, char := range testSuffix {
			assert.True(t, char >= '0' && char <= '9')
		}
	}
}

func TestTransactionID_Comparison(t *testing.T) {
	id1 := NewTransactionID()
	id2 := NewTransactionID()

	// Different instances should have different values (with high probability)
	// Note: There's a small chance they could be the same due to randomness
	// but it's extremely unlikely with proper random generation
	if id1.String() == id2.String() {
		t.Log("Warning: Two generated IDs are the same (very low probability event)")
	}

	// Same string should create equal IDs
	validIDString := "TXN20240729143045123456"
	id3, err := NewTransactionIDFromString(validIDString)
	require.NoError(t, err)

	id4, err := NewTransactionIDFromString(validIDString)
	require.NoError(t, err)

	assert.Equal(t, id3.String(), id4.String())
}

func TestTransactionID_PrefixConsistency(t *testing.T) {
	// Generate multiple IDs and ensure they all have TXN prefix
	for i := 0; i < 100; i++ {
		id := NewTransactionID()
		assert.True(t, strings.HasPrefix(id.String(), "TXN"))
	}
}

func TestTransactionID_LengthConsistency(t *testing.T) {
	// Generate multiple IDs and ensure consistent minimum length
	for i := 0; i < 100; i++ {
		id := NewTransactionID()
		idStr := id.String()

		// Should be exactly 23 characters (TXN + 14 timestamp + 6 suffix)
		assert.Equal(t, 23, len(idStr))
	}
}
