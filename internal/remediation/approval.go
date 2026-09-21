package remediation

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/remediation/isolation"
)

const ApprovalSchemaVersion = "v0alpha1"

type ApprovalErrorCode string

const (
	ApprovalInvalid        ApprovalErrorCode = "APPROVAL_INVALID"
	ApprovalExpired        ApprovalErrorCode = "APPROVAL_EXPIRED"
	ApprovalUsed           ApprovalErrorCode = "APPROVAL_USED"
	ApprovalPlanMismatch   ApprovalErrorCode = "APPROVAL_PLAN_MISMATCH"
	ApprovalSourceMismatch ApprovalErrorCode = "APPROVAL_SOURCE_MISMATCH"
)

type ApprovalError struct {
	Code    ApprovalErrorCode
	Message string
}

func (e *ApprovalError) Error() string { return string(e.Code) + ": " + e.Message }

var approvalUse = struct {
	sync.Mutex
	tokens map[string]bool
}{tokens: map[string]bool{}}

type Approval struct {
	SchemaVersion   string    `json:"schema_version"`
	Token           string    `json:"token"`
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
	canonical, err := isolation.CanonicalRepositoryRoot(context.Background(), sourceRoot)
	if err != nil || plan.RepositoryIdentity.Revision == "" || !expiresAt.After(time.Now().UTC()) {
		return Approval{}, &ApprovalError{Code: ApprovalInvalid, Message: "canonical source root, source commit, and future expiry are required"}
	}
	paths, err := authorizedPaths(plan, canonical)
	if err != nil {
		return Approval{}, err
	}
	tokenBytes := make([]byte, 24)
	if _, err := rand.Read(tokenBytes); err != nil {
		return Approval{}, fmt.Errorf("generate approval token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)
	approvalUse.Lock()
	approvalUse.tokens[token] = false
	approvalUse.Unlock()
	return Approval{SchemaVersion: ApprovalSchemaVersion, Token: token, PlanID: plan.PlanID, PlanDigest: PlanDigest(plan), SourceCommit: plan.RepositoryIdentity.Revision, SourceRoot: canonical, AuthorizedPaths: paths, ExpiresAt: expiresAt.UTC()}, nil
}

func (a Approval) Validate(plan Plan, sourceRoot, sourceCommit string, now time.Time) error {
	if a.SchemaVersion != ApprovalSchemaVersion || a.Token == "" {
		return &ApprovalError{Code: ApprovalInvalid, Message: "approval schema or token is invalid"}
	}
	if !now.Before(a.ExpiresAt) {
		return &ApprovalError{Code: ApprovalExpired, Message: "approval has expired"}
	}
	approvalUse.Lock()
	used, known := approvalUse.tokens[a.Token]
	approvalUse.Unlock()
	if a.Used || (known && used) || !known {
		return &ApprovalError{Code: ApprovalUsed, Message: "approval has already been used or is unknown"}
	}
	if a.PlanID != plan.PlanID || a.PlanDigest != PlanDigest(plan) {
		return &ApprovalError{Code: ApprovalPlanMismatch, Message: "approval does not match plan"}
	}
	canonical, err := isolation.CanonicalRepositoryRoot(context.Background(), sourceRoot)
	if err != nil || a.SourceRoot != canonical || a.SourceCommit != sourceCommit {
		return &ApprovalError{Code: ApprovalSourceMismatch, Message: "approval does not match source identity"}
	}
	paths, err := authorizedPaths(plan, canonical)
	if err != nil {
		return err
	}
	if !equalStrings(a.AuthorizedPaths, paths) {
		return &ApprovalError{Code: ApprovalPlanMismatch, Message: "approval authorized paths do not match plan"}
	}
	return nil
}

func (a *Approval) Consume() error {
	if a == nil || a.Token == "" {
		return &ApprovalError{Code: ApprovalInvalid, Message: "approval token is missing"}
	}
	approvalUse.Lock()
	defer approvalUse.Unlock()
	used, known := approvalUse.tokens[a.Token]
	if !known || used || a.Used {
		return &ApprovalError{Code: ApprovalUsed, Message: "approval has already been used or is unknown"}
	}
	approvalUse.tokens[a.Token] = true
	a.Used = true
	return nil
}

func authorizedPaths(plan Plan, root string) ([]string, error) {
	paths := make([]string, 0, len(plan.AffectedFiles))
	seen := make(map[string]struct{}, len(plan.AffectedFiles))
	for _, file := range plan.AffectedFiles {
		path := filepath.ToSlash(filepath.Clean(filepath.FromSlash(strings.ReplaceAll(file.Path, "\\", "/"))))
		if path == "" || path == "." || filepath.IsAbs(filepath.FromSlash(path)) || path == ".." || strings.HasPrefix(path, "../") {
			return nil, &ApprovalError{Code: ApprovalInvalid, Message: "affected paths must be relative and contained"}
		}
		if _, exists := seen[path]; exists {
			return nil, &ApprovalError{Code: ApprovalInvalid, Message: "affected paths must be unique"}
		}
		seen[path] = struct{}{}
		full := filepath.Join(root, filepath.FromSlash(path))
		if resolved, err := filepath.EvalSymlinks(full); err == nil {
			relative, relErr := filepath.Rel(root, resolved)
			if relErr != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
				return nil, &ApprovalError{Code: ApprovalInvalid, Message: "affected paths must remain inside the repository"}
			}
		} else if !os.IsNotExist(err) {
			return nil, &ApprovalError{Code: ApprovalInvalid, Message: "affected path cannot be resolved safely"}
		}
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths, nil
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func PlanDigest(plan Plan) string {
	canonical := plan
	canonical.PlanID = ""
	canonical.Canonicalize()
	data, _ := json.Marshal(canonical)
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}
