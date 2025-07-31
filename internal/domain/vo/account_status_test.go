package vo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountStatus_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   AccountStatus
		expected bool
	}{
		{
			name:     "Active status is valid",
			status:   AccountStatusActive,
			expected: true,
		},
		{
			name:     "Inactive status is valid",
			status:   AccountStatusInactive,
			expected: true,
		},
		{
			name:     "Suspended status is valid",
			status:   AccountStatusSuspended,
			expected: true,
		},
		{
			name:     "Invalid status",
			status:   AccountStatus("INVALID"),
			expected: false,
		},
		{
			name:     "Empty status",
			status:   AccountStatus(""),
			expected: false,
		},
		{
			name:     "Random string status",
			status:   AccountStatus("RANDOM"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsValid())
		})
	}
}

func TestAccountStatus_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   AccountStatus
		expected bool
	}{
		{
			name:     "Active status",
			status:   AccountStatusActive,
			expected: true,
		},
		{
			name:     "Inactive status",
			status:   AccountStatusInactive,
			expected: false,
		},
		{
			name:     "Suspended status",
			status:   AccountStatusSuspended,
			expected: false,
		},
		{
			name:     "Invalid status",
			status:   AccountStatus("INVALID"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsActive())
		})
	}
}

func TestAccountStatus_IsInactive(t *testing.T) {
	tests := []struct {
		name     string
		status   AccountStatus
		expected bool
	}{
		{
			name:     "Inactive status",
			status:   AccountStatusInactive,
			expected: true,
		},
		{
			name:     "Active status",
			status:   AccountStatusActive,
			expected: false,
		},
		{
			name:     "Suspended status",
			status:   AccountStatusSuspended,
			expected: false,
		},
		{
			name:     "Invalid status",
			status:   AccountStatus("INVALID"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsInactive())
		})
	}
}

func TestAccountStatus_IsSuspended(t *testing.T) {
	tests := []struct {
		name     string
		status   AccountStatus
		expected bool
	}{
		{
			name:     "Suspended status",
			status:   AccountStatusSuspended,
			expected: true,
		},
		{
			name:     "Active status",
			status:   AccountStatusActive,
			expected: false,
		},
		{
			name:     "Inactive status",
			status:   AccountStatusInactive,
			expected: false,
		},
		{
			name:     "Invalid status",
			status:   AccountStatus("INVALID"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsSuspended())
		})
	}
}

func TestAccountStatus_CanTransact(t *testing.T) {
	tests := []struct {
		name     string
		status   AccountStatus
		expected bool
	}{
		{
			name:     "Active account can transact",
			status:   AccountStatusActive,
			expected: true,
		},
		{
			name:     "Inactive account cannot transact",
			status:   AccountStatusInactive,
			expected: false,
		},
		{
			name:     "Suspended account cannot transact",
			status:   AccountStatusSuspended,
			expected: false,
		},
		{
			name:     "Invalid status cannot transact",
			status:   AccountStatus("INVALID"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.CanTransact())
		})
	}
}

func TestAccountStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus AccountStatus
		targetStatus  AccountStatus
		expected      bool
	}{
		// From ACTIVE
		{
			name:          "Active to Inactive",
			currentStatus: AccountStatusActive,
			targetStatus:  AccountStatusInactive,
			expected:      true,
		},
		{
			name:          "Active to Suspended",
			currentStatus: AccountStatusActive,
			targetStatus:  AccountStatusSuspended,
			expected:      true,
		},
		{
			name:          "Active to Active (no change)",
			currentStatus: AccountStatusActive,
			targetStatus:  AccountStatusActive,
			expected:      false,
		},

		// From INACTIVE
		{
			name:          "Inactive to Active",
			currentStatus: AccountStatusInactive,
			targetStatus:  AccountStatusActive,
			expected:      true,
		},
		{
			name:          "Inactive to Suspended",
			currentStatus: AccountStatusInactive,
			targetStatus:  AccountStatusSuspended,
			expected:      true,
		},
		{
			name:          "Inactive to Inactive (no change)",
			currentStatus: AccountStatusInactive,
			targetStatus:  AccountStatusInactive,
			expected:      false,
		},

		// From SUSPENDED
		{
			name:          "Suspended to Active",
			currentStatus: AccountStatusSuspended,
			targetStatus:  AccountStatusActive,
			expected:      true,
		},
		{
			name:          "Suspended to Inactive",
			currentStatus: AccountStatusSuspended,
			targetStatus:  AccountStatusInactive,
			expected:      true,
		},
		{
			name:          "Suspended to Suspended (no change)",
			currentStatus: AccountStatusSuspended,
			targetStatus:  AccountStatusSuspended,
			expected:      false,
		},

		// Invalid statuses
		{
			name:          "Invalid current status",
			currentStatus: AccountStatus("INVALID"),
			targetStatus:  AccountStatusActive,
			expected:      false,
		},
		{
			name:          "Invalid target status",
			currentStatus: AccountStatusActive,
			targetStatus:  AccountStatus("INVALID"),
			expected:      false,
		},
		{
			name:          "Both invalid",
			currentStatus: AccountStatus("INVALID1"),
			targetStatus:  AccountStatus("INVALID2"),
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.currentStatus.CanTransitionTo(tt.targetStatus))
		})
	}
}

func TestAccountStatus_Constants(t *testing.T) {
	// Ensure constants have expected string values
	assert.Equal(t, "ACTIVE", string(AccountStatusActive))
	assert.Equal(t, "INACTIVE", string(AccountStatusInactive))
	assert.Equal(t, "SUSPENDED", string(AccountStatusSuspended))
}

func TestAccountStatus_AllValidStatusesCanTransact(t *testing.T) {
	// Test that we know which statuses can transact
	validStatuses := []AccountStatus{
		AccountStatusActive,
		AccountStatusInactive,
		AccountStatusSuspended,
	}

	expectedCanTransact := map[AccountStatus]bool{
		AccountStatusActive:    true,
		AccountStatusInactive:  false,
		AccountStatusSuspended: false,
	}

	for _, status := range validStatuses {
		t.Run(string(status), func(t *testing.T) {
			assert.Equal(t, expectedCanTransact[status], status.CanTransact())
		})
	}
}

func TestAccountStatus_TransitionMatrix(t *testing.T) {

	// Define expected transitions
	expectedTransitions := map[AccountStatus]map[AccountStatus]bool{
		AccountStatusActive: {
			AccountStatusActive:    false, // Cannot transition to self
			AccountStatusInactive:  true,
			AccountStatusSuspended: true,
		},
		AccountStatusInactive: {
			AccountStatusActive:    true,
			AccountStatusInactive:  false, // Cannot transition to self
			AccountStatusSuspended: true,
		},
		AccountStatusSuspended: {
			AccountStatusActive:    true,
			AccountStatusInactive:  true,
			AccountStatusSuspended: false, // Cannot transition to self
		},
	}

	for fromStatus, transitions := range expectedTransitions {
		for toStatus, expected := range transitions {
			t.Run(string(fromStatus)+"_to_"+string(toStatus), func(t *testing.T) {
				assert.Equal(t, expected, fromStatus.CanTransitionTo(toStatus))
			})
		}
	}
}
