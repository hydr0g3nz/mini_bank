package vo

import (
	"strings"
	"testing"
	"time"

	errs "github.com/hydr0g3nz/mini_bank/internal/domain/error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccountID(t *testing.T) {
	id := NewAccountID()

	assert.NotEmpty(t, id.String())
	assert.Equal(t, 16, len(id.String()))
	assert.True(t, id.IsValid())
	assert.False(t, id.IsEmpty())

	// Check date prefix format (YYYYMMDD)
	datePrefix := id.String()[:8]
	_, err := time.Parse("20060102", datePrefix)
	assert.NoError(t, err)

	// Check that sequence part is numeric
	sequence := id.String()[8:]
	assert.Equal(t, 8, len(sequence))
	for _, char := range sequence {
		assert.True(t, char >= '0' && char <= '9')
	}
}

func TestNewAccountIDFromString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorType   error
	}{
		{
			name:        "Valid account ID",
			input:       "2024072912345678",
			expectError: false,
		},
		{
			name:        "Empty string",
			input:       "",
			expectError: true,
			errorType:   errs.ErrInvalidAccountID,
		},
		{
			name:        "Too short",
			input:       "1234567",
			expectError: true,
			errorType:   errs.ErrInvalidAccountID,
		},
		{
			name:        "Too long",
			input:       "20240729123456789",
			expectError: true,
			errorType:   errs.ErrInvalidAccountID,
		},
		{
			name:        "Invalid date prefix",
			input:       "2024137812345678", // Invalid month
			expectError: true,
			errorType:   errs.ErrInvalidAccountID,
		},
		{
			name:        "Non-numeric characters",
			input:       "2024072912345ABC",
			expectError: true,
			errorType:   errs.ErrInvalidAccountID,
		},
		{
			name:        "Valid edge case - leap year",
			input:       "2024022912345678", // Feb 29 in leap year
			expectError: false,
		},
		{
			name:        "Invalid date - Feb 30",
			input:       "2024023012345678",
			expectError: true,
			errorType:   errs.ErrInvalidAccountID,
		},
		{
			name:        "Current date format",
			input:       time.Now().Format("20060102") + "12345678",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewAccountIDFromString(tt.input)

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

func TestAccountID_String(t *testing.T) {
	validID := "2024072912345678"
	id, err := NewAccountIDFromString(validID)
	require.NoError(t, err)

	assert.Equal(t, validID, id.String())
}

func TestAccountID_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		id       AccountID
		expected bool
	}{
		{
			name:     "Empty account ID",
			id:       AccountID{},
			expected: true,
		},
		{
			name:     "Valid account ID",
			id:       AccountID{value: "2024072912345678"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.id.IsEmpty())
		})
	}
}

func TestAccountID_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		id       AccountID
		expected bool
	}{
		{
			name:     "Valid account ID",
			id:       AccountID{value: "2024072912345678"},
			expected: true,
		},
		{
			name:     "Empty account ID",
			id:       AccountID{},
			expected: false,
		},
		{
			name:     "Invalid format",
			id:       AccountID{value: "invalid"},
			expected: false,
		},
		{
			name:     "Wrong length",
			id:       AccountID{value: "123456789"},
			expected: false,
		},
		{
			name:     "Invalid date",
			id:       AccountID{value: "2024137812345678"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.id.IsValid())
		})
	}
}

func TestValidateAccountID(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "Valid ID",
			input:       "2024072912345678",
			expectError: false,
		},
		{
			name:        "Empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "Wrong length - too short",
			input:       "1234567",
			expectError: true,
		},
		{
			name:        "Wrong length - too long",
			input:       "20240729123456789",
			expectError: true,
		},
		{
			name:        "Non-numeric",
			input:       "2024072912345ABC",
			expectError: true,
		},
		{
			name:        "Invalid date - month 13",
			input:       "2024137812345678",
			expectError: true,
		},
		{
			name:        "Invalid date - day 32",
			input:       "2024073212345678",
			expectError: true,
		},
		{
			name:        "Valid date edge cases",
			input:       "2024123112345678", // Dec 31
			expectError: false,
		},
		{
			name:        "Valid leap year date",
			input:       "2024022912345678", // Feb 29 in leap year 2024
			expectError: false,
		},
		{
			name:        "Short input length check",
			input:       "1234567",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAccountID(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.ErrorIs(t, err, errs.ErrInvalidAccountID)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAccountID_Uniqueness(t *testing.T) {
	// Generate multiple IDs and ensure they are unique
	ids := make(map[string]bool)
	numIDs := 1000

	for i := 0; i < numIDs; i++ {
		id := NewAccountID()
		idStr := id.String()

		// Check that this ID hasn't been generated before
		assert.False(t, ids[idStr], "Duplicate ID generated: %s", idStr)
		ids[idStr] = true

		// Verify format
		assert.Equal(t, 16, len(idStr))
		assert.True(t, id.IsValid())
	}
}

func TestAccountID_DatePrefix(t *testing.T) {
	id := NewAccountID()
	idStr := id.String()

	// Extract date prefix
	datePrefix := idStr[:8]
	today := time.Now().Format("20060102")

	// The date prefix should be today's date
	assert.Equal(t, today, datePrefix)
}

func TestAccountID_SequencePart(t *testing.T) {
	id := NewAccountID()
	idStr := id.String()

	// Extract sequence part
	sequence := idStr[8:]

	// Should be exactly 8 digits
	assert.Equal(t, 8, len(sequence))

	// Should be all numeric
	for _, char := range sequence {
		assert.True(t, char >= '0' && char <= '9', "Non-numeric character in sequence: %c", char)
	}

	// Should be zero-padded (test by generating multiple and checking format)
	for i := 0; i < 10; i++ {
		testID := NewAccountID()
		testSeq := testID.String()[8:]
		assert.Equal(t, 8, len(testSeq))
		assert.True(t, strings.ContainsAny(testSeq, "0123456789"))
	}
}

func TestAccountID_Comparison(t *testing.T) {
	id1 := NewAccountID()
	id2 := NewAccountID()

	// Different instances should have different values (with high probability)
	// Note: There's a small chance they could be the same due to randomness
	// but it's extremely unlikely with proper random generation
	if id1.String() == id2.String() {
		t.Log("Warning: Two generated IDs are the same (very low probability event)")
	}

	// Same string should create equal IDs
	validIDString := "2024072912345678"
	id3, err := NewAccountIDFromString(validIDString)
	require.NoError(t, err)

	id4, err := NewAccountIDFromString(validIDString)
	require.NoError(t, err)

	assert.Equal(t, id3.String(), id4.String())
}
