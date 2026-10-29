package main

import (
	"fmt"
	"testing"
)

func TestFullPayout(t *testing.T) {
	if got := (FullPayout{}).Calculate(2000); got != 2000 {
		t.Errorf("got %.0f, want 2000", got)
	}
}

func TestDeductiblePayout(t *testing.T) {
	if got := (DeductiblePayout{}).Calculate(2000); got != 1500 {
		t.Errorf("got %.0f, want 1500", got)
	}
	if got := (DeductiblePayout{}).Calculate(300); got != 0 {
		t.Errorf("got %.0f, want 0", got)
	}
}

func TestLimitedPayout(t *testing.T) {
	if got := (LimitedPayout{}).Calculate(8000); got != 5000 {
		t.Errorf("got %.0f, want 5000", got)
	}
	if got := (LimitedPayout{}).Calculate(3000); got != 3000 {
		t.Errorf("got %.0f, want 3000", got)
	}
}

func TestPayoutContextSetStrategy(t *testing.T) {
	ctx := NewPayoutContext(FullPayout{})
	if got := ctx.Payout(8000); got != 8000 {
		t.Errorf("full: got %.0f, want 8000", got)
	}
	ctx.SetStrategy(LimitedPayout{})
	if got := ctx.Payout(8000); got != 5000 {
		t.Errorf("limited: got %.0f, want 5000", got)
	}
}

// Parameterized test: the same check for ALL payout strategies.
func TestPayoutBetweenZeroAndDamage(t *testing.T) {
	strategies := []PayoutStrategy{FullPayout{}, DeductiblePayout{}, LimitedPayout{}}
	amounts := []float64{0, 100, 500, 2000, 5000, 20000}
	for _, s := range strategies {
		for _, a := range amounts {
			t.Run(fmt.Sprintf("%T/%.0f", s, a), func(t *testing.T) {
				got := s.Calculate(a)
				if got < 0 || got > a {
					t.Errorf("payout %.0f is outside 0..%.0f", got, a)
				}
			})
		}
	}
}
