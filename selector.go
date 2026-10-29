package main

// Selector chooses a strategy automatically. if/switch is allowed only here.
type Selector struct {
	Risk RiskSource
}

func (s Selector) ChoosePremium(c Customer) PremiumStrategy {
	if c.Years >= 5 {
		return LoyaltyPremium{}
	}
	if c.Age < 25 || c.Age > 60 {
		return AgePremium{}
	}
	if c.Region == "coastal" {
		return RegionPremium{Risk: s.Risk}
	}
	return StandardPremium{}
}
