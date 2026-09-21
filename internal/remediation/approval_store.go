package remediation

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ApprovalStore records consumed approval tokens across CLI processes. The
// marker is created with O_EXCL so concurrent apply attempts cannot reuse one.
type ApprovalStore struct {
	Root string
}

func (s ApprovalStore) Consume(token string) error {
	if token == "" {
		return &ApprovalError{Code: ApprovalInvalid, Message: "approval token is missing"}
	}
	if s.Root == "" {
		return errors.New("approval store root is required")
	}
	if err := os.MkdirAll(s.Root, 0o700); err != nil {
		return fmt.Errorf("create approval store: %w", err)
	}
	if err := os.Chmod(s.Root, 0o700); err != nil {
		return fmt.Errorf("restrict approval store: %w", err)
	}
	digest := sha256.Sum256([]byte(token))
	path := filepath.Join(s.Root, hex.EncodeToString(digest[:])+".used")
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return &ApprovalError{Code: ApprovalUsed, Message: "approval has already been used"}
		}
		return fmt.Errorf("record approval use: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close approval marker: %w", err)
	}
	return nil
}
