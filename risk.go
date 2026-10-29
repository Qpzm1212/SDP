package main

// RiskSource hides where the region risk comes from (Clean Code, Chapter 10).
// Today it is simulated; tomorrow it can be a real API without changing strategies.
type RiskSource interface {
	RegionRisk(region string) float64 // value from 0 to 1
}

// SimulatedRisk is a fake data source with fixed values.
type SimulatedRisk struct{}

func (SimulatedRisk) RegionRisk(region string) float64 {
	risks := map[string]float64{"coastal": 0.8, "city": 0.5, "village": 0.2}
	risk, ok := risks[region]
	if !ok {
		return 0.3
	}
	return risk
}
