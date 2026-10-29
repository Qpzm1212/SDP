package main

import "fmt"

func main() {
	selector := Selector{Risk: SimulatedRisk{}}
	customer := Customer{Name: "Aset", Age: 30, Region: "coastal", Years: 2}

	fmt.Println("=== 1. Same customer, all pricing strategies ===")
	all := []struct {
		name     string
		strategy PremiumStrategy
	}{
		{"standard", StandardPremium{}},
		{"age", AgePremium{}},
		{"loyalty", LoyaltyPremium{}},
		{"region", RegionPremium{Risk: SimulatedRisk{}}},
	}
	for _, s := range all {
		fmt.Printf("%-9s -> %.2f\n", s.name, NewPricingContext(s.strategy).Premium(customer))
	}

	fmt.Println("\n=== 2. Manual vs automatic choice ===")
	ctx := NewPricingContext(StandardPremium{})
	fmt.Printf("manual (standard): %.2f\n", ctx.Premium(customer))
	customers := []Customer{
		{"Young Ali", 20, "city", 0},
		{"Loyal Dana", 45, "city", 8},
		{"Coastal Bek", 35, "coastal", 1},
		{"Regular Sara", 35, "city", 1},
	}
	for _, c := range customers {
		ctx.SetStrategy(selector.ChoosePremium(c))
		fmt.Printf("auto %-13s -> %.2f (%T)\n", c.Name, ctx.Premium(c), selector.ChoosePremium(c))
	}

	fmt.Println("\n=== 3. Switch strategy at runtime ===")
	ctx.SetStrategy(StandardPremium{})
	fmt.Printf("before: %.2f\n", ctx.Premium(customer))
	ctx.SetStrategy(LoyaltyPremium{})
	fmt.Printf("after:  %.2f\n", ctx.Premium(customer))

	fmt.Println("\n=== Payout (old switch version) ===")
	for _, kind := range []string{"full", "deductible", "limited"} {
		fmt.Printf("%-10s for 8000 -> %.2f\n", kind, Payout(kind, 8000))
	}
}
