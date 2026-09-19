package policy

import "testing"

func TestExitCodeMatrix(t *testing.T) {
	cases := []struct {
		outcome Outcome
		want    int
	}{{Pass, 0}, {Warn, 0}, {Block, 4}, {Outcome("unknown"), 3}}
	for _, tc := range cases {
		if got := ExitCode(Decision{Outcome: tc.outcome}); got != tc.want {
			t.Fatalf("ExitCode(%q) = %d, want %d", tc.outcome, got, tc.want)
		}
	}
}
