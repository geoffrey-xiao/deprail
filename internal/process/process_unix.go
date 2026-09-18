//go:build !windows

package process

import (
	"os/exec"
	"syscall"
)

func configureProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func terminateProcessGroup(cmd *exec.Cmd) {
	if cmd.Process != nil {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
}

func exitCode(err error) int {
	if status, ok := err.(*exec.ExitError); ok {
		if wait, ok := status.Sys().(syscall.WaitStatus); ok {
			return wait.ExitStatus()
		}
	}
	return -1
}
