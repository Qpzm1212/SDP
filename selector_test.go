package main

import (
	"fmt"
	"testing"
)

func TestSelectorChoosePremium(t *testing.T) {
	s := Selector{Risk: fakeRisk{value: 0.5}}
	tests := []struct {
		name     string
		customer Customer
		want     string
	}{
		{"loyal customer", Customer{Age: 40, Region: "city", Years: 5}, "main.LoyaltyPremium"},
		{"young customer", Customer{Age: 20, Region: "city", Years: 1}, "main.AgePremium"},
		{"old customer", Customer{Age: 70, Region: "city", Years: 1}, "main.AgePremium"},
		{"coastal customer", Customer{Age: 40, Region: "coastal", Years: 1}, "main.RegionPremium"},
		{"regular customer", Customer{Age: 40, Region: "city", Years: 1}, "main.StandardPremium"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fmt.Sprintf("%T", s.ChoosePremium(tt.customer)); got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestSelectorChoosePayout(t *testing.T) {
	s := Selector{Risk: fakeRisk{value: 0.5}}
	tests := []struct {
		name  string
		claim Claim
		want  string
	}{
		{"huge claim", Claim{Type: "health", Amount: 20000}, "main.LimitedPayout"},
		{"car claim", Claim{Type: "car", Amount: 2000}, "main.DeductiblePayout"},
		{"other claim", Claim{Type: "health", Amount: 2000}, "main.FullPayout"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fmt.Sprintf("%T", s.ChoosePayout(tt.claim)); got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}
