// Assignment 4: Facade over a 3-stage chain. Domain: Insurance Claims.
// A = claim intake, B = policy check, C = payout. Everything lives in this one file.
// A, B and C never mention each other: they only meet inside the facades at the bottom.
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ---------- tiny JSON-file helpers (the "embedded store") ----------

func readJSON[T any](path string) ([]T, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var all []T
	err = json.Unmarshal(data, &all)
	return all, err
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// =====================================================================
// SUBSYSTEM A: claim intake  (Claim, ClaimRules, ClaimStore, ClaimService)
// =====================================================================

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrDocsIncomplete  = errors.New("documents incomplete")
	ErrDuplicateClaim  = errors.New("duplicate claim")
	ErrClaimTransition = errors.New("illegal claim status transition")
	ErrClaimNotFound   = errors.New("claim not found")
)

type ClaimStatus string

const (
	ClaimNew       ClaimStatus = "NEW"
	ClaimAccepted  ClaimStatus = "ACCEPTED"
	ClaimWithdrawn ClaimStatus = "WITHDRAWN"
)

// Claim: entity, holds state and guards its own transitions.
type Claim struct {
	ID       int
	Customer string
	Kind     string
	Status   ClaimStatus
}

func (c *Claim) MoveTo(next ClaimStatus) error {
	ok := (c.Status == ClaimNew && (next == ClaimAccepted || next == ClaimWithdrawn)) ||
		(c.Status == ClaimAccepted && next == ClaimWithdrawn)
	if !ok {
		return fmt.Errorf("%w: %s -> %s", ErrClaimTransition, c.Status, next)
	}
	c.Status = next
	return nil
}

// ClaimRules: pure business logic.
type ClaimRules struct{}

func (ClaimRules) Check(customer, kind string, duplicate bool) error {
	switch {
	case customer == "" || kind == "":
		return ErrInvalidInput
	case customer == "nodocs":
		return fmt.Errorf("%w: customer %s", ErrDocsIncomplete, customer)
	case duplicate:
		return fmt.Errorf("%w: %s/%s", ErrDuplicateClaim, customer, kind)
	}
	return nil
}

// ClaimStore: A's own persistence gateway.
type ClaimStore struct{ path string }

func (s ClaimStore) Load() ([]Claim, error) { return readJSON[Claim](s.path) }
func (s ClaimStore) Save(all []Claim) error { return writeJSON(s.path, all) }

type ClaimService struct {
	rules ClaimRules
	store ClaimStore
}

func NewClaimService(path string) *ClaimService { return &ClaimService{store: ClaimStore{path}} }

func (s *ClaimService) Accept(customer, kind string) (int, error) {
	all, err := s.store.Load()
	if err != nil {
		return 0, err
	}
	dup := false
	for _, x := range all {
		if x.Customer == customer && x.Kind == kind && x.Status == ClaimAccepted {
			dup = true
		}
	}
	if err := s.rules.Check(customer, kind, dup); err != nil {
		return 0, err
	}
	c := Claim{ID: len(all) + 1, Customer: customer, Kind: kind, Status: ClaimNew}
	if err := c.MoveTo(ClaimAccepted); err != nil {
		return 0, err
	}
	return c.ID, s.store.Save(append(all, c))
}

func (s *ClaimService) Withdraw(id int) error {
	all, err := s.store.Load()
	if err != nil {
		return err
	}
	for i := range all {
		if all[i].ID == id {
			if err := all[i].MoveTo(ClaimWithdrawn); err != nil {
				return err
			}
			return s.store.Save(all)
		}
	}
	return ErrClaimNotFound
}

func (s *ClaimService) All() ([]Claim, error) { return s.store.Load() }

// =====================================================================
// SUBSYSTEM B: policy check  (Coverage, CoverageRules, CoverageStore, CoverageService)
// =====================================================================

var (
	ErrPolicyLapsed       = errors.New("policy lapsed")
	ErrNotCovered         = errors.New("not covered by policy")
	ErrCoverageTransition = errors.New("illegal coverage status transition")
	ErrCoverageNotFound   = errors.New("coverage not found")
)

type CoverageStatus string

const (
	CoverageNew       CoverageStatus = "NEW"
	CoverageConfirmed CoverageStatus = "CONFIRMED"
	CoverageRevoked   CoverageStatus = "REVOKED"
)

type Coverage struct {
	ID       int
	Customer string
	Kind     string
	Status   CoverageStatus
}

func (c *Coverage) MoveTo(next CoverageStatus) error {
	ok := (c.Status == CoverageNew && (next == CoverageConfirmed || next == CoverageRevoked)) ||
		(c.Status == CoverageConfirmed && next == CoverageRevoked)
	if !ok {
		return fmt.Errorf("%w: %s -> %s", ErrCoverageTransition, c.Status, next)
	}
	c.Status = next
	return nil
}

type CoverageRules struct{}

func (CoverageRules) Check(customer, kind string) error {
	if customer == "lapsed" {
		return fmt.Errorf("%w: customer %s", ErrPolicyLapsed, customer)
	}
	if kind != "theft" && kind != "fire" && kind != "accident" {
		return fmt.Errorf("%w: %s", ErrNotCovered, kind)
	}
	return nil
}

type CoverageStore struct{ path string }

func (s CoverageStore) Load() ([]Coverage, error) { return readJSON[Coverage](s.path) }
func (s CoverageStore) Save(all []Coverage) error { return writeJSON(s.path, all) }

type CoverageService struct {
	rules CoverageRules
	store CoverageStore
}

func NewCoverageService(path string) *CoverageService {
	return &CoverageService{store: CoverageStore{path}}
}

func (s *CoverageService) Confirm(customer, kind string) (int, error) {
	if err := s.rules.Check(customer, kind); err != nil {
		return 0, err
	}
	all, err := s.store.Load()
	if err != nil {
		return 0, err
	}
	c := Coverage{ID: len(all) + 1, Customer: customer, Kind: kind, Status: CoverageNew}
	if err := c.MoveTo(CoverageConfirmed); err != nil {
		return 0, err
	}
	return c.ID, s.store.Save(append(all, c))
}

func (s *CoverageService) Revoke(id int) error {
	all, err := s.store.Load()
	if err != nil {
		return err
	}
	for i := range all {
		if all[i].ID == id {
			if err := all[i].MoveTo(CoverageRevoked); err != nil {
				return err
			}
			return s.store.Save(all)
		}
	}
	return ErrCoverageNotFound
}

func (s *CoverageService) All() ([]Coverage, error) { return s.store.Load() }

// =====================================================================
// SUBSYSTEM C: payout  (Payment, PaymentRules, PaymentStore, PaymentService)
// =====================================================================

var (
	ErrPayoutBlocked     = errors.New("payout blocked")
	ErrPaymentTransition = errors.New("illegal payment status transition")
	ErrPaymentNotFound   = errors.New("payment not found")
)

type PaymentStatus string

const (
	PaymentNew      PaymentStatus = "NEW"
	PaymentPaid     PaymentStatus = "PAID"
	PaymentReversed PaymentStatus = "REVERSED"
)

type Payment struct {
	ID       int
	Customer string
	Kind     string
	Amount   int
	Status   PaymentStatus
}

func (p *Payment) MoveTo(next PaymentStatus) error {
	ok := (p.Status == PaymentNew && (next == PaymentPaid || next == PaymentReversed)) ||
		(p.Status == PaymentPaid && next == PaymentReversed)
	if !ok {
		return fmt.Errorf("%w: %s -> %s", ErrPaymentTransition, p.Status, next)
	}
	p.Status = next
	return nil
}

type PaymentRules struct{}

func (PaymentRules) Amount(kind string) int {
	return map[string]int{"theft": 500, "fire": 2000, "accident": 800}[kind]
}

func (PaymentRules) Check(customer string) error {
	if customer == "frozen" {
		return fmt.Errorf("%w: bank account of %s is frozen", ErrPayoutBlocked, customer)
	}
	return nil
}

type PaymentStore struct{ path string }

func (s PaymentStore) Load() ([]Payment, error) { return readJSON[Payment](s.path) }
func (s PaymentStore) Save(all []Payment) error { return writeJSON(s.path, all) }

type PaymentService struct {
	rules PaymentRules
	store PaymentStore
}

func NewPaymentService(path string) *PaymentService {
	return &PaymentService{store: PaymentStore{path}}
}

func (s *PaymentService) Pay(customer, kind string) (int, error) {
	if err := s.rules.Check(customer); err != nil {
		return 0, err
	}
	all, err := s.store.Load()
	if err != nil {
		return 0, err
	}
	p := Payment{ID: len(all) + 1, Customer: customer, Kind: kind, Amount: s.rules.Amount(kind), Status: PaymentNew}
	if err := p.MoveTo(PaymentPaid); err != nil {
		return 0, err
	}
	return p.ID, s.store.Save(append(all, p))
}

func (s *PaymentService) Reverse(id int) error {
	all, err := s.store.Load()
	if err != nil {
		return err
	}
	for i := range all {
		if all[i].ID == id {
			if err := all[i].MoveTo(PaymentReversed); err != nil {
				return err
			}
			return s.store.Save(all)
		}
	}
	return ErrPaymentNotFound
}

func (s *PaymentService) All() ([]Payment, error) { return s.store.Load() }

// =====================================================================
// FACADES (the only place where A, B and C meet)
// =====================================================================

var (
	ErrChainNotFound      = errors.New("chain not found")
	ErrAlreadyCompensated = errors.New("chain already compensated")
)

type ChainStatus string

const (
	PendingB    ChainStatus = "PENDING_B"   // A committed, B not yet
	PendingC    ChainStatus = "PENDING_C"   // A and B committed, C not yet
	Completed   ChainStatus = "COMPLETED"   // A, B, C all committed
	Compensated ChainStatus = "COMPENSATED" // everything undone
)

// Chain = one claim attempt, saved to disk after every step.
type Chain struct {
	ID                             int
	Customer, Kind                 string
	ClaimID, CoverageID, PaymentID int
	Status                         ChainStatus
}

// deps: what both facades share (the 3 subsystems + the chain log file).
type deps struct {
	a       *ClaimService
	b       *CoverageService
	c       *PaymentService
	logPath string
}

func newDeps(dir string) (*deps, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &deps{
		a:       NewClaimService(filepath.Join(dir, "claims.json")),
		b:       NewCoverageService(filepath.Join(dir, "coverages.json")),
		c:       NewPaymentService(filepath.Join(dir, "payments.json")),
		logPath: filepath.Join(dir, "chains.json"),
	}, nil
}

func (d *deps) chains() ([]Chain, error) { return readJSON[Chain](d.logPath) }

// save inserts the chain or replaces the one with the same ID.
func (d *deps) save(c Chain) error {
	all, err := d.chains()
	if err != nil {
		return err
	}
	found := false
	for i := range all {
		if all[i].ID == c.ID {
			all[i], found = c, true
		}
	}
	if !found {
		all = append(all, c)
	}
	return writeJSON(d.logPath, all)
}

// compensate undoes a chain. WHAT it undoes depends on the status read from disk:
//
//	COMPLETED -> reverse C, revoke B, withdraw A
//	PENDING_C -> revoke B, withdraw A      (C failed)
//	PENDING_B -> withdraw A                (B failed)
func (d *deps) compensate(c *Chain) error {
	if c.Status == Compensated {
		return ErrAlreadyCompensated
	}
	if c.Status == Completed {
		if err := d.c.Reverse(c.PaymentID); err != nil {
			return err
		}
	}
	if c.Status == Completed || c.Status == PendingC {
		if err := d.b.Revoke(c.CoverageID); err != nil {
			return err
		}
	}
	if err := d.a.Withdraw(c.ClaimID); err != nil {
		return err
	}
	c.Status = Compensated
	return d.save(*c)
}

// ---- Facade 1: normal client workflow ----

type ClaimFacade struct{ *deps }

func NewClaimFacade(dir string) (*ClaimFacade, error) {
	d, err := newDeps(dir)
	if err != nil {
		return nil, err
	}
	return &ClaimFacade{d}, nil
}

// Submit runs A -> B -> C. If B or C fails, it compensates and returns the error.
func (f *ClaimFacade) Submit(customer, kind string) (Chain, error) {
	all, err := f.chains()
	if err != nil {
		return Chain{}, err
	}
	chain := Chain{ID: len(all) + 1, Customer: customer, Kind: kind}

	// Step A. Fails => nothing committed, nothing to undo.
	chain.ClaimID, err = f.a.Accept(customer, kind)
	if err != nil {
		return Chain{}, err
	}
	chain.Status = PendingB
	if err := f.save(chain); err != nil {
		return chain, err
	}

	// Step B. Fails => undo A only.
	chain.CoverageID, err = f.b.Confirm(customer, kind)
	if err != nil {
		return f.fail(chain, err)
	}
	chain.Status = PendingC
	if err := f.save(chain); err != nil {
		return chain, err
	}

	// Step C. Fails => undo B, then A.
	chain.PaymentID, err = f.c.Pay(customer, kind)
	if err != nil {
		return f.fail(chain, err)
	}
	chain.Status = Completed
	return chain, f.save(chain)
}

func (f *ClaimFacade) fail(chain Chain, cause error) (Chain, error) {
	if err := f.compensate(&chain); err != nil {
		return chain, err
	}
	return chain, fmt.Errorf("claim rolled back: %w", cause)
}

// Recover (call at startup): compensates every chain left half-done on disk.
func (f *ClaimFacade) Recover() (int, error) {
	all, err := f.chains()
	if err != nil {
		return 0, err
	}
	fixed := 0
	for _, c := range all {
		if c.Status == PendingB || c.Status == PendingC {
			if err := f.compensate(&c); err != nil {
				return fixed, err
			}
			fixed++
		}
	}
	return fixed, nil
}

// ---- Facade 2: admin / inspection ----

type AdminFacade struct{ *deps }

func NewAdminFacade(dir string) (*AdminFacade, error) {
	d, err := newDeps(dir)
	if err != nil {
		return nil, err
	}
	return &AdminFacade{d}, nil
}

type Report struct {
	Chains    []Chain
	Claims    []Claim
	Coverages []Coverage
	Payments  []Payment
}

func (f *AdminFacade) Inspect() (Report, error) {
	var r Report
	var err error
	if r.Chains, err = f.chains(); err != nil {
		return r, err
	}
	if r.Claims, err = f.a.All(); err != nil {
		return r, err
	}
	if r.Coverages, err = f.b.All(); err != nil {
		return r, err
	}
	r.Payments, err = f.c.All()
	return r, err
}

// Cancel manually compensates one chain (works on COMPLETED too).
func (f *AdminFacade) Cancel(chainID int) error {
	all, err := f.chains()
	if err != nil {
		return err
	}
	for _, c := range all {
		if c.ID == chainID {
			return f.compensate(&c)
		}
	}
	return ErrChainNotFound
}

// =====================================================================
// DEMO APPLICATION (console menu)
// =====================================================================

const menu = `
1) Submit a claim
2) Inspect saved state
3) Demo: B fails (lapsed + theft)  => withdraw A
4) Demo: C fails (frozen + theft)  => revoke B, then withdraw A
5) Cancel a chain manually
6) Reset all data
0) Exit
Customers: alice, nodocs, lapsed, frozen   Kinds: theft, fire, accident, flood`

func run(dataDir string, r io.Reader, w io.Writer) {
	cf, err := NewClaimFacade(dataDir)
	if err != nil {
		fmt.Fprintln(w, "error:", err)
		return
	}
	admin, err := NewAdminFacade(dataDir)
	if err != nil {
		fmt.Fprintln(w, "error:", err)
		return
	}
	if n, err := cf.Recover(); err != nil {
		fmt.Fprintln(w, "recover error:", err)
	} else if n > 0 {
		fmt.Fprintf(w, "Recovered %d unfinished chain(s) from disk and compensated them\n", n)
	}

	sc := bufio.NewScanner(r)
	ask := func(prompt string) string {
		fmt.Fprint(w, prompt)
		sc.Scan()
		return strings.TrimSpace(sc.Text())
	}
	submit := func(customer, kind string) {
		chain, err := cf.Submit(customer, kind)
		if err != nil {
			fmt.Fprintln(w, "RESULT: FAILED ->", err)
			if chain.ID != 0 {
				fmt.Fprintf(w, "chain #%d is now %s\n", chain.ID, chain.Status)
			} else {
				fmt.Fprintln(w, "nothing was committed")
			}
			return
		}
		fmt.Fprintf(w, "RESULT: OK, chain #%d is %s\n", chain.ID, chain.Status)
	}

	for {
		fmt.Fprintln(w, menu)
		fmt.Fprint(w, "> ")
		if !sc.Scan() {
			return
		}
		switch strings.TrimSpace(sc.Text()) {
		case "1":
			customer := ask("customer: ")
			submit(customer, ask("kind: "))
		case "2":
			rep, err := admin.Inspect()
			if err != nil {
				fmt.Fprintln(w, "error:", err)
				break
			}
			for _, x := range rep.Chains {
				fmt.Fprintf(w, "CHAIN    %+v\n", x)
			}
			for _, x := range rep.Claims {
				fmt.Fprintf(w, "CLAIM    %+v\n", x)
			}
			for _, x := range rep.Coverages {
				fmt.Fprintf(w, "COVERAGE %+v\n", x)
			}
			for _, x := range rep.Payments {
				fmt.Fprintf(w, "PAYMENT  %+v\n", x)
			}
		case "3":
			submit("lapsed", "theft")
		case "4":
			submit("frozen", "theft")
		case "5":
			id, err := strconv.Atoi(ask("chain id: "))
			if err != nil {
				fmt.Fprintln(w, "not a number")
				break
			}
			if err := admin.Cancel(id); err != nil {
				fmt.Fprintln(w, "error:", err)
			} else {
				fmt.Fprintln(w, "chain compensated")
			}
		case "6":
			os.RemoveAll(dataDir)
			os.MkdirAll(dataDir, 0o755) // keep the folder so the app keeps working
			fmt.Fprintln(w, "data deleted")
		case "0":
			return
		}
	}
}

func main() { run("data", os.Stdin, os.Stdout) }
