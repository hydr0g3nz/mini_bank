package vo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionStatus_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   TransactionStatus
		expected bool
	}{
		{
			name:     "Pending status is valid",
			status:   TransactionStatusPending,
			expected: true,
		},
		{
			name:     "Completed status is valid",
			status:   TransactionStatusCompleted,
			expected: true,
		},
		{
			name:     "Failed status is valid",
			status:   TransactionStatusFailed,
			expected: true,
		},
		{
			name:     "Cancelled status is valid",
			status:   TransactionStatusCancelled,
			expected: true,
		},
		{
			name:     "Invalid status",
			status:   TransactionStatus("INVALID"),
			expected: false,
		},
		{
			name:     "Empty status",
			status:   TransactionStatus(""),
			expected: false,
		},
		{
			name:     "Random string status",
			status:   TransactionStatus("RANDOM"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsValid())
		})
	}
}

func TestTransactionStatus_IsPending(t *testing.T) {
	tests := []struct {
		name     string
		status   TransactionStatus
		expected bool
	}{
		{
			name:     "Pending status",
			status:   TransactionStatusPending,
			expected: true,
		},
		{
			name:     "Completed status",
			status:   TransactionStatusCompleted,
			expected: false,
		},
		{
			name:     "Failed status",
			status:   TransactionStatusFailed,
			expected: false,
		},
		{
			name:     "Cancelled status",
			status:   TransactionStatusCancelled,
			expected: false,
		},
		{
			name:     "Invalid status",
			status:   TransactionStatus("INVALID"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsPending())
		})
	}
}

func TestTransactionStatus_IsCompleted(t *testing.T) {
	tests := []struct {
		name     string
		status   TransactionStatus
		expected bool
	}{
		{
			name:     "Completed status",
			status:   TransactionStatusCompleted,
			expected: true,
		},
		{
			name:     "Pending status",
			status:   TransactionStatusPending,
			expected: false,
		},
		{
			name:     "Failed status",
			status:   TransactionStatusFailed,
			expected: false,
		},
		{
			name:     "Cancelled status",
			status:   TransactionStatusCancelled,
			expected: false,
		},
		{
			name:     "Invalid status",
			status:   TransactionStatus("INVALID"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsCompleted())
		})
	}
}

func TestTransactionStatus_IsFailed(t *testing.T) {
	tests := []struct {
		name     string
		status   TransactionStatus
		expected bool
	}{
		{
			name:     "Failed status",
			status:   TransactionStatusFailed,
			expected: true,
		},
		{
			name:     "Pending status",
			status:   TransactionStatusPending,
			expected: false,
		},
		{
			name:     "Completed status",
			status:   TransactionStatusCompleted,
			expected: false,
		},
		{
			name:     "Cancelled status",
			status:   TransactionStatusCancelled,
			expected: false,
		},
		{
			name:     "Invalid status",
			status:   TransactionStatus("INVALID"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsFailed())
		})
	}
}

func TestTransactionStatus_IsCancelled(t *testing.T) {
	tests := []struct {
		name     string
		status   TransactionStatus
		expected bool
	}{
		{
			name:     "Cancelled status",
			status:   TransactionStatusCancelled,
			expected: true,
		},
		{
			name:     "Pending status",
			status:   TransactionStatusPending,
			expected: false,
		},
		{
			name:     "Completed status",
			status:   TransactionStatusCompleted,
			expected: false,
		},
		{
			name:     "Failed status",
			status:   TransactionStatusFailed,
			expected: false,
		},
		{
			name:     "Invalid status",
			status:   TransactionStatus("INVALID"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsCancelled())
		})
	}
}

func TestTransactionStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus TransactionStatus
		targetStatus  TransactionStatus
		expected      bool
	}{
		// From PENDING
		{
			name:          "Pending to Completed",
			currentStatus: TransactionStatusPending,
			targetStatus:  TransactionStatusCompleted,
			expected:      true,
		},
		{
			name:          "Pending to Failed",
			currentStatus: TransactionStatusPending,
			targetStatus:  TransactionStatusFailed,
			expected:      true,
		},
		{
			name:          "Pending to Cancelled",
			currentStatus: TransactionStatusPending,
			targetStatus:  TransactionStatusCancelled,
			expected:      true,
		},
		{
			name:          "Pending to Pending (no change)",
			currentStatus: TransactionStatusPending,
			targetStatus:  TransactionStatusPending,
			expected:      false,
		},

		// From COMPLETED
		{
			name:          "Completed to Pending (invalid)",
			currentStatus: TransactionStatusCompleted,
			targetStatus:  TransactionStatusPending,
			expected:      false,
		},
		{
			name:          "Completed to Failed (invalid)",
			currentStatus: TransactionStatusCompleted,
			targetStatus:  TransactionStatusFailed,
			expected:      false,
		},
		{
			name:          "Completed to Cancelled (invalid)",
			currentStatus: TransactionStatusCompleted,
			targetStatus:  TransactionStatusCancelled,
			expected:      false,
		},
		{
			name:          "Completed to Completed (no change)",
			currentStatus: TransactionStatusCompleted,
			targetStatus:  TransactionStatusCompleted,
			expected:      false,
		},

		// From FAILED
		{
			name:          "Failed to Cancelled",
			currentStatus: TransactionStatusFailed,
			targetStatus:  TransactionStatusCancelled,
			expected:      true,
		},
		{
			name:          "Failed to Pending (invalid)",
			currentStatus: TransactionStatusFailed,
			targetStatus:  TransactionStatusPending,
			expected:      false,
		},
		{
			name:          "Failed to Completed (invalid)",
			currentStatus: TransactionStatusFailed,
			targetStatus:  TransactionStatusCompleted,
			expected:      false,
		},
		{
			name:          "Failed to Failed (no change)",
			currentStatus: TransactionStatusFailed,
			targetStatus:  TransactionStatusFailed,
			expected:      false,
		},

		// From CANCELLED
		{
			name:          "Cancelled to Pending (invalid)",
			currentStatus: TransactionStatusCancelled,
			targetStatus:  TransactionStatusPending,
			expected:      false,
		},
		{
			name:          "Cancelled to Completed (invalid)",
			currentStatus: TransactionStatusCancelled,
			targetStatus:  TransactionStatusCompleted,
			expected:      false,
		},
		{
			name:          "Cancelled to Failed (invalid)",
			currentStatus: TransactionStatusCancelled,
			targetStatus:  TransactionStatusFailed,
			expected:      false,
		},
		{
			name:          "Cancelled to Cancelled (no change)",
			currentStatus: TransactionStatusCancelled,
			targetStatus:  TransactionStatusCancelled,
			expected:      false,
		},

		// Invalid statuses
		{
			name:          "Invalid current status",
			currentStatus: TransactionStatus("INVALID"),
			targetStatus:  TransactionStatusCompleted,
			expected:      false,
		},
		{
			name:          "Invalid target status",
			currentStatus: TransactionStatusPending,
			targetStatus:  TransactionStatus("INVALID"),
			expected:      false,
		},
		{
			name:          "Both invalid",
			currentStatus: TransactionStatus("INVALID1"),
			targetStatus:  TransactionStatus("INVALID2"),
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.currentStatus.CanTransitionTo(tt.targetStatus))
		})
	}
}

func TestTransactionStatus_Constants(t *testing.T) {
	// Ensure constants have expected string values
	assert.Equal(t, "PENDING", string(TransactionStatusPending))
	assert.Equal(t, "COMPLETED", string(TransactionStatusCompleted))
	assert.Equal(t, "FAILED", string(TransactionStatusFailed))
	assert.Equal(t, "CANCELLED", string(TransactionStatusCancelled))
}

func TestTransactionStatus_TransitionMatrix(t *testing.T) {

	// Define expected transitions
	expectedTransitions := map[TransactionStatus]map[TransactionStatus]bool{
		TransactionStatusPending: {
			TransactionStatusPending:   false, // Cannot transition to self
			TransactionStatusCompleted: true,
			TransactionStatusFailed:    true,
			TransactionStatusCancelled: true,
		},
		TransactionStatusCompleted: {
			TransactionStatusPending:   false,
			TransactionStatusCompleted: false, // Cannot transition to self
			TransactionStatusFailed:    false,
			TransactionStatusCancelled: false,
		},
		TransactionStatusFailed: {
			TransactionStatusPending:   false,
			TransactionStatusCompleted: false,
			TransactionStatusFailed:    false, // Cannot transition to self
			TransactionStatusCancelled: true,
		},
		TransactionStatusCancelled: {
			TransactionStatusPending:   false,
			TransactionStatusCompleted: false,
			TransactionStatusFailed:    false,
			TransactionStatusCancelled: false, // Cannot transition to self
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

func TestTransactionStatus_FinalStates(t *testing.T) {
	// Test that final states cannot transition to any other state
	finalStates := []TransactionStatus{
		TransactionStatusCompleted,
		TransactionStatusCancelled,
	}

	allStates := []TransactionStatus{
		TransactionStatusPending,
		TransactionStatusCompleted,
		TransactionStatusFailed,
		TransactionStatusCancelled,
	}

	for _, finalState := range finalStates {
		for _, targetState := range allStates {
			t.Run(string(finalState)+"_cannot_transition_to_"+string(targetState), func(t *testing.T) {
				// Final states should not be able to transition to any state
				// (including themselves, since that's not a real transition)
				assert.False(t, finalState.CanTransitionTo(targetState))
			})
		}
	}
}

func TestTransactionStatus_PendingTransitions(t *testing.T) {
	// Test that pending can transition to all other states
	nonPendingStates := []TransactionStatus{
		TransactionStatusCompleted,
		TransactionStatusFailed,
		TransactionStatusCancelled,
	}

	for _, targetState := range nonPendingStates {
		t.Run("Pending_can_transition_to_"+string(targetState), func(t *testing.T) {
			assert.True(t, TransactionStatusPending.CanTransitionTo(targetState))
		})
	}

	// But cannot transition to itself
	assert.False(t, TransactionStatusPending.CanTransitionTo(TransactionStatusPending))
}

func TestTransactionStatus_FailedTransitions(t *testing.T) {
	// Test that failed can only transition to cancelled
	assert.True(t, TransactionStatusFailed.CanTransitionTo(TransactionStatusCancelled))

	// Cannot transition to other states
	assert.False(t, TransactionStatusFailed.CanTransitionTo(TransactionStatusPending))
	assert.False(t, TransactionStatusFailed.CanTransitionTo(TransactionStatusCompleted))
	assert.False(t, TransactionStatusFailed.CanTransitionTo(TransactionStatusFailed))
}
