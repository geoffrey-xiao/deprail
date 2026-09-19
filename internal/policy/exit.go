package policy

// ExitCode maps a policy decision to the stable CLI contract.
func ExitCode(decision Decision) int {
	switch decision.Outcome {
	case Pass, Warn:
		return 0
	case Block:
		return 4
	case Unusable:
		return 3
	default:
		return 3
	}
}
