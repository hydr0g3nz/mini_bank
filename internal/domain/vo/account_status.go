package vo

type AccountStatus string

const (
	AccountStatusActive    AccountStatus = "ACTIVE"
	AccountStatusInactive  AccountStatus = "INACTIVE"
	AccountStatusSuspended AccountStatus = "SUSPENDED"
)

// IsValid checks if account status is valid
func (s AccountStatus) IsValid() bool {
	switch s {
	case AccountStatusActive, AccountStatusInactive, AccountStatusSuspended:
		return true
	default:
		return false
	}
}

// IsActive checks if account is active
func (s AccountStatus) IsActive() bool {
	return s == AccountStatusActive
}

// IsInactive checks if account is inactive
func (s AccountStatus) IsInactive() bool {
	return s == AccountStatusInactive
}

// IsSuspended checks if account is suspended
func (s AccountStatus) IsSuspended() bool {
	return s == AccountStatusSuspended
}

// CanTransact checks if account can perform transactions
func (s AccountStatus) CanTransact() bool {
	return s == AccountStatusActive
}

// CanTransitionTo checks if current status can transition to target status
func (s AccountStatus) CanTransitionTo(target AccountStatus) bool {
	switch s {
	case AccountStatusActive:
		return target == AccountStatusInactive || target == AccountStatusSuspended
	case AccountStatusInactive:
		return target == AccountStatusActive || target == AccountStatusSuspended
	case AccountStatusSuspended:
		return target == AccountStatusActive || target == AccountStatusInactive
	default:
		return false
	}
}
