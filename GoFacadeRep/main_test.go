package main

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

// ---------- tiny domain-specific testing language ----------
// Tests read like: s.customerClaims("alice","theft") ... s.thenChainIs(1, Completed)

type scenario struct {
	t     *testing.T
	dir   string
	cf    *ClaimFacade
	admin *AdminFacade
}

// Every test gets its own empty folder => independent + repeatable.
func newScenario(t *testing.T) *scenario {
	dir := t.TempDir()
	cf, err := NewClaimFacade(dir)
	if err != nil {
		t.Fatal(err)
	}
	admin, err := NewAdminFacade(dir)
	if err != nil {
		t.Fatal(err)
	}
	return &scenario{t, dir, cf, admin}
}

func (s *scenario) customerClaims(customer, kind string) (Chain, error) {
	return s.cf.Submit(customer, kind)
}

func (s *scenario) report() Report {
	r, err := s.admin.Inspect()
	if err != nil {
		s.t.Fatal(err)
	}
	return r
}

func (s *scenario) thenChainIs(id int, want ChainStatus) {
	s.t.Helper()
	for _, c := range s.report().Chains {
		if c.ID == id && c.Status == want {
			return
		}
	}
	s.t.Fatalf("chain %d is not %s: %+v", id, want, s.report().Chains)
}

func (s *scenario) thenClaimIs(want ClaimStatus) {
	s.t.Helper()
	if x := s.report().Claims; len(x) != 1 || x[0].Status != want {
		s.t.Fatalf("claims = %+v, want one %s", x, want)
	}
}

func (s *scenario) thenCoverageIs(want CoverageStatus) {
	s.t.Helper()
	if x := s.report().Coverages; len(x) != 1 || x[0].Status != want {
		s.t.Fatalf("coverages = %+v, want one %s", x, want)
	}
}

func (s *scenario) thenPaymentIs(want PaymentStatus) {
	s.t.Helper()
	if x := s.report().Payments; len(x) != 1 || x[0].Status != want {
		s.t.Fatalf("payments = %+v, want one %s", x, want)
	}
}

func (s *scenario) thenNothingIsSaved() {
	s.t.Helper()
	r := s.report()
	if len(r.Chains)+len(r.Claims)+len(r.Coverages)+len(r.Payments) != 0 {
		s.t.Fatalf("expected empty store, got %+v", r)
	}
}

// ---------- facade tests ----------

func TestHappyPath(t *testing.T) {
	s := newScenario(t)
	if chain, err := s.customerClaims("alice", "theft"); err != nil || chain.Status != Completed {
		t.Fatalf("got %+v, %v", chain, err)
	}
	s.thenChainIs(1, Completed)
	s.thenClaimIs(ClaimAccepted)
	s.thenCoverageIs(CoverageConfirmed)
	s.thenPaymentIs(PaymentPaid)
}

// Depth 1: B fails after A committed -> only A is undone.
func TestBFailsAfterA_CompensatesAOnly(t *testing.T) {
	s := newScenario(t)
	if _, err := s.customerClaims("lapsed", "theft"); !errors.Is(err, ErrPolicyLapsed) {
		t.Fatalf("want ErrPolicyLapsed, got %v", err)
	}
	s.thenChainIs(1, Compensated)
	s.thenClaimIs(ClaimWithdrawn)
	if r := s.report(); len(r.Coverages) != 0 || len(r.Payments) != 0 {
		t.Fatalf("B and C must be untouched: %+v", r)
	}
}

// Depth 2: C fails after A and B committed -> B is undone, then A.
func TestCFailsAfterAB_CompensatesBThenA(t *testing.T) {
	s := newScenario(t)
	if _, err := s.customerClaims("frozen", "theft"); !errors.Is(err, ErrPayoutBlocked) {
		t.Fatalf("want ErrPayoutBlocked, got %v", err)
	}
	s.thenChainIs(1, Compensated)
	s.thenCoverageIs(CoverageRevoked)
	s.thenClaimIs(ClaimWithdrawn)
	if r := s.report(); len(r.Payments) != 0 {
		t.Fatalf("C must be untouched: %+v", r)
	}
}

// Extra failures at the very first step: nothing is committed.
func TestFirstStepFailures_NothingCommitted(t *testing.T) {
	s := newScenario(t)
	if _, err := s.customerClaims("nodocs", "theft"); !errors.Is(err, ErrDocsIncomplete) {
		t.Fatalf("got %v", err)
	}
	if _, err := s.customerClaims("", "theft"); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("got %v", err)
	}
	s.thenNothingIsSaved()
}

func TestDuplicateAndNotCovered(t *testing.T) {
	s := newScenario(t)
	s.customerClaims("alice", "theft")
	if _, err := s.customerClaims("alice", "theft"); !errors.Is(err, ErrDuplicateClaim) {
		t.Fatalf("got %v", err)
	}
	if _, err := s.customerClaims("bob", "flood"); !errors.Is(err, ErrNotCovered) {
		t.Fatalf("got %v", err)
	}
}

// Restart test: "crash" after A and B, new objects on the same files, Recover() reads the disk.
func TestRestartMidChain_RecoverCompensatesFromDisk(t *testing.T) {
	s := newScenario(t)
	aid, _ := s.cf.a.Accept("alice", "theft")
	bid, _ := s.cf.b.Confirm("alice", "theft")
	s.cf.save(Chain{ID: 1, Customer: "alice", Kind: "theft", ClaimID: aid, CoverageID: bid, Status: PendingC})

	restarted, _ := NewClaimFacade(s.dir)
	if n, err := restarted.Recover(); err != nil || n != 1 {
		t.Fatalf("recovered %d, err %v", n, err)
	}
	s.thenChainIs(1, Compensated)
	s.thenCoverageIs(CoverageRevoked)
	s.thenClaimIs(ClaimWithdrawn)
	if n, _ := restarted.Recover(); n != 0 {
		t.Fatalf("second recover must do nothing, did %d", n)
	}
}

func TestRestartMidChain_PendingBOnlyWithdrawsA(t *testing.T) {
	s := newScenario(t)
	aid, _ := s.cf.a.Accept("alice", "theft")
	s.cf.save(Chain{ID: 1, Customer: "alice", Kind: "theft", ClaimID: aid, Status: PendingB})

	restarted, _ := NewClaimFacade(s.dir)
	if n, err := restarted.Recover(); err != nil || n != 1 {
		t.Fatalf("recovered %d, err %v", n, err)
	}
	s.thenClaimIs(ClaimWithdrawn)
}

func TestAdminCancel(t *testing.T) {
	s := newScenario(t)
	s.customerClaims("alice", "theft")
	if err := s.admin.Cancel(1); err != nil {
		t.Fatal(err)
	}
	s.thenChainIs(1, Compensated)
	s.thenPaymentIs(PaymentReversed)
	s.thenCoverageIs(CoverageRevoked)
	s.thenClaimIs(ClaimWithdrawn)
	if err := s.admin.Cancel(1); !errors.Is(err, ErrAlreadyCompensated) {
		t.Fatalf("got %v", err)
	}
	if err := s.admin.Cancel(42); !errors.Is(err, ErrChainNotFound) {
		t.Fatalf("got %v", err)
	}
}

// ---------- each subsystem used alone, without the facade ----------

func TestSubsystemsStandalone(t *testing.T) {
	dir := t.TempDir()

	a := NewClaimService(filepath.Join(dir, "a.json"))
	id, err := a.Accept("alice", "fire")
	if err != nil || a.Withdraw(id) != nil {
		t.Fatalf("A failed: %v", err)
	}
	if !errors.Is(a.Withdraw(id), ErrClaimTransition) || !errors.Is(a.Withdraw(99), ErrClaimNotFound) {
		t.Fatal("A must guard its transitions")
	}

	b := NewCoverageService(filepath.Join(dir, "b.json"))
	id, err = b.Confirm("alice", "fire")
	if err != nil || b.Revoke(id) != nil {
		t.Fatalf("B failed: %v", err)
	}
	if !errors.Is(b.Revoke(id), ErrCoverageTransition) || !errors.Is(b.Revoke(99), ErrCoverageNotFound) {
		t.Fatal("B must guard its transitions")
	}

	c := NewPaymentService(filepath.Join(dir, "c.json"))
	id, err = c.Pay("alice", "fire")
	if err != nil || c.Reverse(id) != nil {
		t.Fatalf("C failed: %v", err)
	}
	if !errors.Is(c.Reverse(id), ErrPaymentTransition) || !errors.Is(c.Reverse(99), ErrPaymentNotFound) {
		t.Fatal("C must guard its transitions")
	}
	if all, _ := c.All(); all[0].Amount != 2000 {
		t.Fatalf("fire payout should be 2000, got %+v", all)
	}
}

// ---------- the demo app, driven by a scripted "user" ----------

func TestMenuDemo(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "data")
	input := "1\nalice\ntheft\n3\n4\n2\n5\n1\n5\nx\n5\n99\n6\n0\n"
	var out bytes.Buffer
	run(dir, strings.NewReader(input), &out)
	for _, want := range []string{
		"RESULT: OK", "policy lapsed", "payout blocked", "COMPENSATED",
		"chain compensated", "not a number", "chain not found", "data deleted",
	} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("output is missing %q:\n%s", want, out.String())
		}
	}
}
