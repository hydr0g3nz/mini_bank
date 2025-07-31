// domain/vo/money.go
package vo

import (
	"errors"

	"github.com/shopspring/decimal"
)

// Money represents monetary value using Shopspring decimal for precision
type Money struct {
	amount decimal.Decimal
}

// NewMoney creates a new Money instance from decimal
func NewMoney(amount decimal.Decimal) Money {
	return Money{
		amount: amount,
	}
}

// NewMoneyFromString creates Money from string representation
func NewMoneyFromString(amount string) (Money, error) {
	dec, err := decimal.NewFromString(amount)
	if err != nil {
		return Money{}, err
	}
	return NewMoney(dec), nil
}

// NewMoneyFromFloat creates Money from float64
func NewMoneyFromFloat(amount float64) Money {
	return NewMoney(decimal.NewFromFloat(amount))
}

// NewMoneyFromInt creates Money from int64
func NewMoneyFromInt(amount int64) Money {
	return NewMoney(decimal.NewFromInt(amount))
}

// Zero returns a Money instance with zero value
func ZeroMoney() Money {
	return NewMoney(decimal.Zero)
}

// Amount returns the decimal amount
func (m Money) Amount() decimal.Decimal {
	return m.amount
}

// IsZero checks if amount is zero
func (m Money) IsZero() bool {
	return m.amount.IsZero()
}

// IsPositive checks if amount is positive
func (m Money) IsPositive() bool {
	return m.amount.IsPositive()
}

// IsNegative checks if amount is negative
func (m Money) IsNegative() bool {
	return m.amount.IsNegative()
}

// Add adds two Money values
func (m Money) Add(other Money) (Money, error) {
	return Money{
		amount: m.amount.Add(other.amount),
	}, nil
}

// Subtract subtracts two Money values
func (m Money) Subtract(other Money) (Money, error) {
	return Money{
		amount: m.amount.Sub(other.amount),
	}, nil
}

// Multiply multiplies Money by a decimal factor
func (m Money) Multiply(factor decimal.Decimal) Money {
	return Money{
		amount: m.amount.Mul(factor),
	}
}

// MultiplyFloat multiplies Money by a float64 factor
func (m Money) MultiplyFloat(factor float64) Money {
	return m.Multiply(decimal.NewFromFloat(factor))
}

// Divide divides Money by a decimal divisor
func (m Money) Divide(divisor decimal.Decimal) (Money, error) {
	if divisor.IsZero() {
		return Money{}, errors.New("cannot divide by zero")
	}
	return Money{
		amount: m.amount.Div(divisor),
	}, nil
}

// DivideFloat divides Money by a float64 divisor
func (m Money) DivideFloat(divisor float64) (Money, error) {
	return m.Divide(decimal.NewFromFloat(divisor))
}

// Abs returns the absolute value of Money
func (m Money) Abs() Money {
	return Money{
		amount: m.amount.Abs(),
	}
}

// Equal checks if two Money values are equal
func (m Money) Equal(other Money) bool {
	return m.amount.Equal(other.amount)
}

// GreaterThan checks if this Money is greater than other
func (m Money) GreaterThan(other Money) bool {
	return m.amount.GreaterThan(other.amount)
}

// GreaterThanOrEqual checks if this Money is greater than or equal to other
func (m Money) GreaterThanOrEqual(other Money) bool {
	return m.amount.GreaterThanOrEqual(other.amount)
}

// LessThan checks if this Money is less than other
func (m Money) LessThan(other Money) bool {
	return m.amount.LessThan(other.amount)
}

// LessThanOrEqual checks if this Money is less than or equal to other
func (m Money) LessThanOrEqual(other Money) bool {
	return m.amount.LessThanOrEqual(other.amount)
}

// Round rounds the Money to the specified number of decimal places
func (m Money) Round(places int32) Money {
	return Money{
		amount: m.amount.Round(places),
	}
}

// RoundBank rounds the Money using banker's rounding
func (m Money) RoundBank(places int32) Money {
	return Money{
		amount: m.amount.RoundBank(places),
	}
}

// Truncate truncates the Money to the specified number of decimal places
func (m Money) Truncate(places int32) Money {
	return Money{
		amount: m.amount.Truncate(places),
	}
}

// String returns string representation
func (m Money) String() string {
	return m.amount.String()
}

// StringFixed returns string representation with fixed decimal places
func (m Money) StringFixed(places int32) string {
	return m.amount.StringFixed(places)
}

// Float64 returns float64 representation (use with caution for precision)
func (m Money) Float64() float64 {
	f, _ := m.amount.Float64()
	return f
}

// InexactFloat64 returns float64 representation (may lose precision)
func (m Money) InexactFloat64() float64 {
	return m.amount.InexactFloat64()
}

// IntPart returns the integer part of the Money
func (m Money) IntPart() int64 {
	return m.amount.IntPart()
}

// Copy returns a copy of the Money value
func (m Money) Copy() Money {
	return Money{
		amount: m.amount.Copy(),
	}
}
