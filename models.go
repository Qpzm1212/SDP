package main

// Customer is a person who buys an insurance policy.
type Customer struct {
	Name   string
	Age    int
	Region string
	Years  int // years with our company
}

// Claim is a request to receive money after an insured event.
type Claim struct {
	Type   string // "car", "health", ...
	Amount float64
}
