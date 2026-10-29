package main

// PayoutStrategy is one way to calculate how much money the customer receives.
type PayoutStrategy interface {
	Calculate(amount float64) float64
}

// FullPayout pays the whole damage.
type FullPayout struct{}

func (FullPayout) Calculate(amount float64) float64 {
	return amount
}

// DeductiblePayout subtracts a fixed deductible from the damage.
type DeductiblePayout struct{}

func (DeductiblePayout) Calculate(amount float64) float64 {
	result := amount - 500
	if result < 0 {
		return 0
	}
	return result
}

// LimitedPayout never pays more than the policy limit.
type LimitedPayout struct{}

func (LimitedPayout) Calculate(amount float64) float64 {
	if amount > 5000 {
		return 5000
	}
	return amount
}

// PayoutContext holds the current strategy and just calls it.
type PayoutContext struct {
	strategy PayoutStrategy
}

func NewPayoutContext(s PayoutStrategy) *PayoutContext {
	return &PayoutContext{strategy: s}
}

func (p *PayoutContext) SetStrategy(s PayoutStrategy) {
	p.strategy = s
}

func (p *PayoutContext) Payout(amount float64) float64 {
	return p.strategy.Calculate(amount)
}

// payoutByName lets old code (and the old tests) choose a strategy by name.
var payoutByName = map[string]PayoutStrategy{
	"full":       FullPayout{},
	"deductible": DeductiblePayout{},
	"limited":    LimitedPayout{},
}

// Payout keeps the same signature as before the refactoring.
func Payout(kind string, amount float64) float64 {
	return NewPayoutContext(payoutByName[kind]).Payout(amount)
}
