//go:build windows

package process

import "os/exec"

func configureProcessGroup(cmd *exec.Cmd) {}

func terminateProcessGroup(cmd *exec.Cmd) {
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
}

func exitCode(err error) int {
	if status, ok := err.(*exec.ExitError); ok {
		return status.ExitCode()
	}
	return -1
}
