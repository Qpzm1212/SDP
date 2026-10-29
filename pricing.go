package main

const basePremium = 1000.0

// PremiumStrategy is one way to calculate the yearly insurance price.
type PremiumStrategy interface {
	Calculate(c Customer) float64
}

// StandardPremium: everybody pays the base price.
type StandardPremium struct{}

func (StandardPremium) Calculate(c Customer) float64 {
	return basePremium
}

// AgePremium: young and old customers pay more.
type AgePremium struct{}

func (AgePremium) Calculate(c Customer) float64 {
	if c.Age < 25 {
		return basePremium * 1.5
	}
	if c.Age > 60 {
		return basePremium * 1.2
	}
	return basePremium
}

// LoyaltyPremium: every year with us gives a discount (max 10 years).
type LoyaltyPremium struct{}

func (LoyaltyPremium) Calculate(c Customer) float64 {
	years := c.Years
	if years > 10 {
		years = 10
	}
	return basePremium - float64(years)*30
}

// RegionPremium: price depends on the risk of the region (external data).
type RegionPremium struct {
	Risk RiskSource
}

func (r RegionPremium) Calculate(c Customer) float64 {
	return basePremium + r.Risk.RegionRisk(c.Region)*500
}

// PricingContext holds the current strategy and just calls it.
type PricingContext struct {
	strategy PremiumStrategy
}

func NewPricingContext(s PremiumStrategy) *PricingContext {
	return &PricingContext{strategy: s}
}

func (p *PricingContext) SetStrategy(s PremiumStrategy) {
	p.strategy = s
}

func (p *PricingContext) Premium(c Customer) float64 {
	return p.strategy.Calculate(c)
}
