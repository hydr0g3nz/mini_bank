package vo

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeDebit    TransactionType = "DEBIT"
	TransactionTypeCredit   TransactionType = "CREDIT"
	TransactionTypeTransfer TransactionType = "TRANSFER"
)

// IsValid checks if transaction type is valid
func (t TransactionType) IsValid() bool {
	switch t {
	case TransactionTypeDebit, TransactionTypeCredit, TransactionTypeTransfer:
		return true
	default:
		return false
	}
}

// IsDebit checks if transaction type is debit
func (t TransactionType) IsDebit() bool {
	return t == TransactionTypeDebit
}

// IsCredit checks if transaction type is credit
func (t TransactionType) IsCredit() bool {
	return t == TransactionTypeCredit
}

// IsTransfer checks if transaction type is transfer
func (t TransactionType) IsTransfer() bool {
	return t == TransactionTypeTransfer
}
