package vo

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
	TransactionStatusCancelled TransactionStatus = "CANCELLED"
)

// IsValid checks if transaction status is valid
func (s TransactionStatus) IsValid() bool {
	switch s {
	case TransactionStatusPending, TransactionStatusCompleted,
		TransactionStatusFailed, TransactionStatusCancelled:
		return true
	default:
		return false
	}
}

// IsPending checks if status is pending
func (s TransactionStatus) IsPending() bool {
	return s == TransactionStatusPending
}

// IsCompleted checks if status is completed
func (s TransactionStatus) IsCompleted() bool {
	return s == TransactionStatusCompleted
}

// IsFailed checks if status is failed
func (s TransactionStatus) IsFailed() bool {
	return s == TransactionStatusFailed
}

// IsCancelled checks if status is cancelled
func (s TransactionStatus) IsCancelled() bool {
	return s == TransactionStatusCancelled
}

// CanTransitionTo checks if current status can transition to target status
func (s TransactionStatus) CanTransitionTo(target TransactionStatus) bool {
	switch s {
	case TransactionStatusPending:
		return target == TransactionStatusCompleted ||
			target == TransactionStatusFailed ||
			target == TransactionStatusCancelled
	case TransactionStatusCompleted:
		return false // Completed transactions cannot be changed
	case TransactionStatusFailed:
		return target == TransactionStatusCancelled
	case TransactionStatusCancelled:
		return false // Cancelled transactions cannot be changed
	default:
		return false
	}
}
