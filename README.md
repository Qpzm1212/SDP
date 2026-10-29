# Insurance: Strategy Pattern (Go)

Domain: **Insurance**. Two strategy families:

1. **Pricing**: how the yearly premium is calculated.
2. **Payout**: how much money the customer receives after a claim.

## Run

```
go run .
go test ./... -v
```

## Classes and responsibilities

| Class / file | Responsibility |
|---|---|
| `Customer` (models.go) | Holds the personal data used for the premium calculation. |
| `Claim` (models.go) | Holds the claim type plus the damage amount used for the payout. |
| `PremiumStrategy` (pricing.go) | Interface with one method that returns the premium for a customer. |
| `StandardPremium` | Returns the base price for every customer. |
| `AgePremium` | Raises the price for young customers, elderly customers. |
| `LoyaltyPremium` | Lowers the price for each year the customer stays with us. |
| `RegionPremium` | Adds a surcharge taken from the region risk data source. |
| `PricingContext` | Stores the current premium strategy, delegates the calculation to it. |
| `RiskSource` (risk.go) | Interface that hides where region risk data comes from. |
| `SimulatedRisk` | Returns fixed region risk values instead of calling a real service. |
| `PayoutStrategy` (payout.go) | Interface with one method that returns the payout for a damage amount. |
| `FullPayout` | Pays the whole damage amount. |
| `DeductiblePayout` | Subtracts a fixed deductible from the damage amount. |
| `LimitedPayout` | Never pays more than the policy limit. |
| `PayoutContext` | Stores the current payout strategy, delegates the calculation to it. |
| `Selector` (selector.go) | Picks a premium strategy for a customer, a payout strategy for a claim. |

## Why Strategy, not State or Template Method

**Not State.** In State the object changes its own behavior when its internal condition changes,
and the states usually know the next state (draft -> active -> closed). Our algorithms do not know
each other and do not move the context anywhere. The client (or the `Selector`) picks one from the
outside, depending on the input data. A customer does not "become" a loyalty price; we choose a
price formula for this one calculation.

**Not Template Method.** Template Method fixes the skeleton of an algorithm in a base class and
lets subclasses override some steps. It uses inheritance and the choice is made at compile time.
Here every algorithm is a full replacement, there is no shared skeleton, and Go has no class
inheritance anyway. We use composition: the context holds an interface and we can call
`SetStrategy()` while the program is running.

**Why it is Strategy.** One interface, several interchangeable implementations, a context that only
delegates, and the choice of implementation is made by the client at runtime. A new algorithm is a
new type that implements the interface. `PricingContext`, `PayoutContext` and the existing
strategies stay unchanged (Open/Closed Principle). The only place that knows the rules of choice is
`Selector`, so `if` is used only there.

## Isolation from change

`RegionPremium` needs the region risk. It does not call any API. It depends on the `RiskSource`
interface. The program uses `SimulatedRisk`; the tests use `fakeRisk` with a fixed value.

## Refactoring history

The payout family was first written the bad way (one function, `switch` on a string) and then
refactored to Strategy. The file `payout_behavior_test.go` is the same in both commits and passes
in both.

- Before (switch): `9875c2f`  ->  `https://github.com/<your-login>/<your-repo>/commit/9875c2f`
- After (Strategy): `36edb19`  ->  `https://github.com/<your-login>/<your-repo>/commit/36edb19`

Local check: `git show 9875c2f` and `git show 36edb19`.

## Tests

- One test per strategy (`pricing_test.go`, `payout_test.go`).
- One test per `Selector` rule (`selector_test.go`).
- Parameterized (table-driven) tests that run the same check on all strategies of a family:
  `TestPremiumNeverNegative`, `TestPayoutBetweenZeroAndDamage`.
