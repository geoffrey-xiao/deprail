package remediation

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Approval struct {
	PlanID          string    `json:"plan_id"`
	PlanDigest      string    `json:"plan_digest"`
	SourceCommit    string    `json:"source_commit"`
	SourceRoot      string    `json:"source_root"`
	AuthorizedPaths []string  `json:"authorized_paths"`
	ExpiresAt       time.Time `json:"expires_at"`
	Used            bool      `json:"used"`
}

func NewApproval(plan Plan, sourceRoot string, expiresAt time.Time) (Approval, error) {
	if err := plan.Validate(); err != nil {
		return Approval{}, fmt.Errorf("validate plan: %w", err)
	}
	if sourceRoot == "" || plan.RepositoryIdentity.Revision == "" || !expiresAt.After(time.Now().UTC()) {
		return Approval{}, errors.New("source root, source commit, and future expiry are required")
	}
	paths := make([]string, 0, len(plan.AffectedFiles))
	for _, file := range plan.AffectedFiles {
		if file.Path == "" || filepath.IsAbs(file.Path) || filepath.Clean(file.Path) == "." || strings.HasPrefix(filepath.Clean(file.Path), ".."+string(filepath.Separator)) {
			return Approval{}, errors.New("affected paths must be relative and contained")
		}
		paths = append(paths, filepath.ToSlash(filepath.Clean(file.Path)))
	}
	sort.Strings(paths)
	return Approval{PlanID: plan.PlanID, PlanDigest: PlanDigest(plan), SourceCommit: plan.RepositoryIdentity.Revision, SourceRoot: sourceRoot, AuthorizedPaths: paths, ExpiresAt: expiresAt.UTC()}, nil
}

func (a Approval) Validate(plan Plan, sourceRoot, sourceCommit string, now time.Time) error {
	if a.Used {
		return errors.New("approval has already been used")
	}
	if !now.Before(a.ExpiresAt) {
		return errors.New("approval has expired")
	}
	if a.PlanID != plan.PlanID || a.PlanDigest != PlanDigest(plan) {
		return errors.New("approval does not match plan")
	}
	if a.SourceRoot != sourceRoot || a.SourceCommit != sourceCommit {
		return errors.New("approval does not match source identity")
	}
	return nil
}

func (a *Approval) Consume() error {
	if a == nil || a.Used {
		return errors.New("approval has already been used")
	}
	a.Used = true
	return nil
}

func PlanDigest(plan Plan) string {
	canonical := plan
	canonical.PlanID = ""
	canonical.Canonicalize()
	data, _ := json.Marshal(canonical)
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}
