package main

import "testing"

// These tests describe the behavior of Payout.
// They must pass BEFORE and AFTER the refactoring (same file in both commits).
func TestPayoutBehavior(t *testing.T) {
	tests := []struct {
		kind   string
		amount float64
		want   float64
	}{
		{"full", 2000, 2000},
		{"deductible", 2000, 1500},
		{"deductible", 300, 0},
		{"limited", 8000, 5000},
		{"limited", 3000, 3000},
	}
	for _, tt := range tests {
		got := Payout(tt.kind, tt.amount)
		if got != tt.want {
			t.Errorf("Payout(%s, %.0f) = %.0f, want %.0f", tt.kind, tt.amount, got, tt.want)
		}
	}
}
