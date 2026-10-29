package main

import (
	"fmt"
	"testing"
)

// fakeRisk is a fake data source: always returns the same value.
type fakeRisk struct{ value float64 }

func (f fakeRisk) RegionRisk(region string) float64 { return f.value }

func TestStandardPremium(t *testing.T) {
	if got := (StandardPremium{}).Calculate(Customer{}); got != 1000 {
		t.Errorf("got %.2f, want 1000", got)
	}
}

func TestAgePremium(t *testing.T) {
	tests := []struct {
		age  int
		want float64
	}{{20, 1500}, {40, 1000}, {70, 1200}}
	for _, tt := range tests {
		if got := (AgePremium{}).Calculate(Customer{Age: tt.age}); got != tt.want {
			t.Errorf("age %d: got %.2f, want %.2f", tt.age, got, tt.want)
		}
	}
}

func TestLoyaltyPremium(t *testing.T) {
	tests := []struct {
		years int
		want  float64
	}{{0, 1000}, {5, 850}, {20, 700}} // 20 years is limited to 10
	for _, tt := range tests {
		if got := (LoyaltyPremium{}).Calculate(Customer{Years: tt.years}); got != tt.want {
			t.Errorf("years %d: got %.2f, want %.2f", tt.years, got, tt.want)
		}
	}
}

func TestRegionPremiumUsesFakeData(t *testing.T) {
	s := RegionPremium{Risk: fakeRisk{value: 0.5}}
	if got := s.Calculate(Customer{Region: "anything"}); got != 1250 {
		t.Errorf("got %.2f, want 1250", got)
	}
}

func TestPricingContextSetStrategy(t *testing.T) {
	ctx := NewPricingContext(StandardPremium{})
	c := Customer{Years: 10}
	if got := ctx.Premium(c); got != 1000 {
		t.Errorf("standard: got %.2f, want 1000", got)
	}
	ctx.SetStrategy(LoyaltyPremium{})
	if got := ctx.Premium(c); got != 700 {
		t.Errorf("loyalty: got %.2f, want 700", got)
	}
}

// Parameterized test: the same check for ALL pricing strategies.
func TestPremiumNeverNegative(t *testing.T) {
	strategies := []PremiumStrategy{
		StandardPremium{},
		AgePremium{},
		LoyaltyPremium{},
		RegionPremium{Risk: fakeRisk{value: 1}},
	}
	customers := []Customer{
		{Age: 0, Years: 0}, {Age: 20, Years: 3}, {Age: 80, Years: 50}, {Age: 40, Years: 10},
	}
	for _, s := range strategies {
		for _, c := range customers {
			name := fmt.Sprintf("%T/age%d/years%d", s, c.Age, c.Years)
			t.Run(name, func(t *testing.T) {
				if got := s.Calculate(c); got < 0 {
					t.Errorf("premium is negative: %.2f", got)
				}
			})
		}
	}
}
