package remediation

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestPersistedApprovalValidatesWithoutProcessLocalToken(t *testing.T) {
	plan, root := approvalTestPlan(t)
	approval, err := NewApproval(plan, root, time.Now().UTC().Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	validated := approval
	approvalUse.Lock()
	delete(approvalUse.tokens, validated.Token)
	approvalUse.Unlock()
	if err := validated.ValidatePersisted(plan, root, "commit-1", time.Now().UTC()); err != nil {
		t.Fatalf("persisted validation: %v", err)
	}
}

func TestApprovalStoreConsumesExactlyOnceConcurrently(t *testing.T) {
	store := ApprovalStore{Root: t.TempDir()}
	const attempts = 12
	results := make(chan error, attempts)
	var wait sync.WaitGroup
	for i := 0; i < attempts; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			results <- store.Consume("token-concurrent")
		}()
	}
	wait.Wait()
	close(results)
	var successes, used int
	for err := range results {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, os.ErrExist):
			used++
		default:
			var approvalErr *ApprovalError
			if !errors.As(err, &approvalErr) || approvalErr.Code != ApprovalUsed {
				t.Fatalf("unexpected consume error: %v", err)
			}
			used++
		}
	}
	if successes != 1 || used != attempts-1 {
		t.Fatalf("successes=%d used=%d", successes, used)
	}
	entries, err := os.ReadDir(store.Root)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("marker count=%d, want 1", len(entries))
	}
	info, err := os.Stat(filepath.Join(store.Root, entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("marker permissions=%o, want 600", info.Mode().Perm())
	}
}
