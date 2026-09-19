package policy

import (
	"errors"
	"strings"
	"time"
)

type Exception struct {
	SchemaVersion   string    `json:"schema_version"`
	ID              string    `json:"id"`
	FindingKey      string    `json:"finding_key"`
	Scope           string    `json:"scope"`
	Rationale       string    `json:"rationale"`
	Approver        string    `json:"approver"`
	CreatedAt       time.Time `json:"created_at"`
	ExpiresAt       time.Time `json:"expires_at"`
	ReviewCondition string    `json:"review_condition"`
}

func (e Exception) Validate() error {
	if e.SchemaVersion != "v1alpha" || strings.TrimSpace(e.ID) == "" || strings.TrimSpace(e.FindingKey) == "" || strings.TrimSpace(e.Scope) == "" || strings.TrimSpace(e.Rationale) == "" || strings.TrimSpace(e.Approver) == "" || strings.TrimSpace(e.ReviewCondition) == "" {
		return errors.New("exception required fields are missing")
	}
	if e.CreatedAt.IsZero() || e.ExpiresAt.IsZero() || !e.ExpiresAt.After(e.CreatedAt) {
		return errors.New("exception timestamps are invalid")
	}
	return nil
}

func (e Exception) ActiveAt(now time.Time) bool {
	return e.Validate() == nil && now.Before(e.ExpiresAt)
}
func (e Exception) Matches(findingKey, scope string, now time.Time) bool {
	return e.ActiveAt(now) && e.FindingKey == findingKey && e.Scope == scope
}
