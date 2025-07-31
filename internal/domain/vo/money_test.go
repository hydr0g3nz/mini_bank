package vo

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMoney(t *testing.T) {
	amount := decimal.NewFromFloat(100.50)
	money := NewMoney(amount)

	assert.True(t, money.Amount().Equal(amount))
	assert.Equal(t, "100.5", money.String())
}

func TestNewMoneyFromString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expected    string
	}{
		{
			name:        "Valid decimal string",
			input:       "100.50",
			expectError: false,
			expected:    "100.5",
		},
		{
			name:        "Integer string",
			input:       "100",
			expectError: false,
			expected:    "100",
		},
		{
			name:        "Zero value",
			input:       "0",
			expectError: false,
			expected:    "0",
		},
		{
			name:        "Negative value",
			input:       "-50.25",
			expectError: false,
			expected:    "-50.25",
		},
		{
			name:        "Invalid string",
			input:       "invalid",
			expectError: true,
		},
		{
			name:        "Empty string",
			input:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := NewMoneyFromString(tt.input)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, money.String())
			}
		})
	}
}

func TestNewMoneyFromFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{
			name:     "Positive float",
			input:    100.50,
			expected: "100.5",
		},
		{
			name:     "Zero",
			input:    0.0,
			expected: "0",
		},
		{
			name:     "Negative float",
			input:    -50.25,
			expected: "-50.25",
		},
		{
			name:     "Large number",
			input:    999999.99,
			expected: "999999.99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money := NewMoneyFromFloat(tt.input)
			assert.Equal(t, tt.expected, money.String())
		})
	}
}

func TestNewMoneyFromInt(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{
			name:     "Positive integer",
			input:    100,
			expected: "100",
		},
		{
			name:     "Zero",
			input:    0,
			expected: "0",
		},
		{
			name:     "Negative integer",
			input:    -50,
			expected: "-50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money := NewMoneyFromInt(tt.input)
			assert.Equal(t, tt.expected, money.String())
		})
	}
}

func TestZeroMoney(t *testing.T) {
	money := ZeroMoney()
	assert.True(t, money.IsZero())
	assert.Equal(t, "0", money.String())
}

func TestMoney_IsZero(t *testing.T) {
	tests := []struct {
		name     string
		money    Money
		expected bool
	}{
		{
			name:     "Zero money",
			money:    ZeroMoney(),
			expected: true,
		},
		{
			name:     "Positive money",
			money:    NewMoneyFromFloat(100.0),
			expected: false,
		},
		{
			name:     "Negative money",
			money:    NewMoneyFromFloat(-100.0),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.money.IsZero())
		})
	}
}

func TestMoney_IsPositive(t *testing.T) {
	tests := []struct {
		name     string
		money    Money
		expected bool
	}{
		{
			name:     "Positive money",
			money:    NewMoneyFromFloat(100.0),
			expected: true,
		},
		{
			name:     "Zero money",
			money:    ZeroMoney(),
			expected: false,
		},
		{
			name:     "Negative money",
			money:    NewMoneyFromFloat(-100.0),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.money.IsPositive())
		})
	}
}

func TestMoney_IsNegative(t *testing.T) {
	tests := []struct {
		name     string
		money    Money
		expected bool
	}{
		{
			name:     "Negative money",
			money:    NewMoneyFromFloat(-100.0),
			expected: true,
		},
		{
			name:     "Zero money",
			money:    ZeroMoney(),
			expected: false,
		},
		{
			name:     "Positive money",
			money:    NewMoneyFromFloat(100.0),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.money.IsNegative())
		})
	}
}

func TestMoney_Add(t *testing.T) {
	tests := []struct {
		name     string
		money1   Money
		money2   Money
		expected string
	}{
		{
			name:     "Add positive numbers",
			money1:   NewMoneyFromFloat(100.50),
			money2:   NewMoneyFromFloat(50.25),
			expected: "150.75",
		},
		{
			name:     "Add zero",
			money1:   NewMoneyFromFloat(100.0),
			money2:   ZeroMoney(),
			expected: "100",
		},
		{
			name:     "Add negative number",
			money1:   NewMoneyFromFloat(100.0),
			money2:   NewMoneyFromFloat(-50.0),
			expected: "50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.money1.Add(tt.money2)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result.String())
		})
	}
}

func TestMoney_Subtract(t *testing.T) {
	tests := []struct {
		name     string
		money1   Money
		money2   Money
		expected string
	}{
		{
			name:     "Subtract positive numbers",
			money1:   NewMoneyFromFloat(100.50),
			money2:   NewMoneyFromFloat(50.25),
			expected: "50.25",
		},
		{
			name:     "Subtract zero",
			money1:   NewMoneyFromFloat(100.0),
			money2:   ZeroMoney(),
			expected: "100",
		},
		{
			name:     "Subtract negative number",
			money1:   NewMoneyFromFloat(100.0),
			money2:   NewMoneyFromFloat(-50.0),
			expected: "150",
		},
		{
			name:     "Result in negative",
			money1:   NewMoneyFromFloat(50.0),
			money2:   NewMoneyFromFloat(100.0),
			expected: "-50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.money1.Subtract(tt.money2)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result.String())
		})
	}
}

func TestMoney_Multiply(t *testing.T) {
	tests := []struct {
		name     string
		money    Money
		factor   decimal.Decimal
		expected string
	}{
		{
			name:     "Multiply by 2",
			money:    NewMoneyFromFloat(100.0),
			factor:   decimal.NewFromFloat(2.0),
			expected: "200",
		},
		{
			name:     "Multiply by 0.5",
			money:    NewMoneyFromFloat(100.0),
			factor:   decimal.NewFromFloat(0.5),
			expected: "50",
		},
		{
			name:     "Multiply by zero",
			money:    NewMoneyFromFloat(100.0),
			factor:   decimal.Zero,
			expected: "0",
		},
		{
			name:     "Multiply by negative",
			money:    NewMoneyFromFloat(100.0),
			factor:   decimal.NewFromFloat(-1.0),
			expected: "-100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.money.Multiply(tt.factor)
			assert.Equal(t, tt.expected, result.String())
		})
	}
}

func TestMoney_MultiplyFloat(t *testing.T) {
	tests := []struct {
		name     string
		money    Money
		factor   float64
		expected string
	}{
		{
			name:     "Multiply by 1.5",
			money:    NewMoneyFromFloat(100.0),
			factor:   1.5,
			expected: "150",
		},
		{
			name:     "Multiply by 0",
			money:    NewMoneyFromFloat(100.0),
			factor:   0,
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.money.MultiplyFloat(tt.factor)
			assert.Equal(t, tt.expected, result.String())
		})
	}
}

func TestMoney_Divide(t *testing.T) {
	tests := []struct {
		name        string
		money       Money
		divisor     decimal.Decimal
		expected    string
		expectError bool
	}{
		{
			name:     "Divide by 2",
			money:    NewMoneyFromFloat(100.0),
			divisor:  decimal.NewFromFloat(2.0),
			expected: "50",
		},
		{
			name:     "Divide by 4",
			money:    NewMoneyFromFloat(100.0),
			divisor:  decimal.NewFromFloat(4.0),
			expected: "25",
		},
		{
			name:        "Divide by zero",
			money:       NewMoneyFromFloat(100.0),
			divisor:     decimal.Zero,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.money.Divide(tt.divisor)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result.String())
			}
		})
	}
}

func TestMoney_DivideFloat(t *testing.T) {
	tests := []struct {
		name        string
		money       Money
		divisor     float64
		expected    string
		expectError bool
	}{
		{
			name:     "Divide by 2.5",
			money:    NewMoneyFromFloat(100.0),
			divisor:  2.5,
			expected: "40",
		},
		{
			name:        "Divide by zero",
			money:       NewMoneyFromFloat(100.0),
			divisor:     0.0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.money.DivideFloat(tt.divisor)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result.String())
			}
		})
	}
}

func TestMoney_Abs(t *testing.T) {
	tests := []struct {
		name     string
		money    Money
		expected string
	}{
		{
			name:     "Positive number",
			money:    NewMoneyFromFloat(100.0),
			expected: "100",
		},
		{
			name:     "Negative number",
			money:    NewMoneyFromFloat(-100.0),
			expected: "100",
		},
		{
			name:     "Zero",
			money:    ZeroMoney(),
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.money.Abs()
			assert.Equal(t, tt.expected, result.String())
		})
	}
}

func TestMoney_Comparisons(t *testing.T) {
	money100 := NewMoneyFromFloat(100.0)
	money50 := NewMoneyFromFloat(50.0)
	money100_2 := NewMoneyFromFloat(100.0)

	t.Run("Equal", func(t *testing.T) {
		assert.True(t, money100.Equal(money100_2))
		assert.False(t, money100.Equal(money50))
	})

	t.Run("GreaterThan", func(t *testing.T) {
		assert.True(t, money100.GreaterThan(money50))
		assert.False(t, money50.GreaterThan(money100))
		assert.False(t, money100.GreaterThan(money100_2))
	})

	t.Run("GreaterThanOrEqual", func(t *testing.T) {
		assert.True(t, money100.GreaterThanOrEqual(money50))
		assert.True(t, money100.GreaterThanOrEqual(money100_2))
		assert.False(t, money50.GreaterThanOrEqual(money100))
	})

	t.Run("LessThan", func(t *testing.T) {
		assert.True(t, money50.LessThan(money100))
		assert.False(t, money100.LessThan(money50))
		assert.False(t, money100.LessThan(money100_2))
	})

	t.Run("LessThanOrEqual", func(t *testing.T) {
		assert.True(t, money50.LessThanOrEqual(money100))
		assert.True(t, money100.LessThanOrEqual(money100_2))
		assert.False(t, money100.LessThanOrEqual(money50))
	})
}

func TestMoney_Round(t *testing.T) {
	tests := []struct {
		name     string
		money    Money
		places   int32
		expected string
	}{
		{
			name:     "Round to 2 places",
			money:    newMoneyFromStringMustValue("100.456"),
			places:   2,
			expected: "100.46",
		},
		{
			name:     "Round to 0 places",
			money:    newMoneyFromStringMustValue("100.456"),
			places:   0,
			expected: "100",
		},
		{
			name:     "Round down",
			money:    newMoneyFromStringMustValue("100.444"),
			places:   2,
			expected: "100.44",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.money.Round(tt.places)
			assert.Equal(t, tt.expected, result.String())
		})
	}
}

func TestMoney_StringFixed(t *testing.T) {
	money := NewMoneyFromFloat(100.5)

	tests := []struct {
		name     string
		places   int32
		expected string
	}{
		{
			name:     "2 decimal places",
			places:   2,
			expected: "100.50",
		},
		{
			name:     "0 decimal places",
			places:   0,
			expected: "101", // Note: this might round
		},
		{
			name:     "4 decimal places",
			places:   4,
			expected: "100.5000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := money.StringFixed(tt.places)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMoney_Copy(t *testing.T) {
	original := NewMoneyFromFloat(100.0)
	copy := original.Copy()

	assert.True(t, original.Equal(copy))

	// Modify original (through operations that return new Money)
	modified, _ := original.Add(NewMoneyFromFloat(50.0))

	// Copy should remain unchanged
	assert.True(t, copy.Equal(NewMoneyFromFloat(100.0)))
	assert.False(t, copy.Equal(modified))
}

// Helper method for tests
func newMoneyFromStringMustValue(s string) Money {
	money, err := NewMoneyFromString(s)
	if err != nil {
		panic(err)
	}
	return money
}
