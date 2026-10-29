package main

// Payout calculates how much money the customer receives.
// BAD VERSION: all variants are in one function with switch (before refactoring).
func Payout(kind string, amount float64) float64 {
	switch kind {
	case "full":
		return amount
	case "deductible":
		result := amount - 500
		if result < 0 {
			return 0
		}
		return result
	case "limited":
		if amount > 5000 {
			return 5000
		}
		return amount
	}
	return 0
}
