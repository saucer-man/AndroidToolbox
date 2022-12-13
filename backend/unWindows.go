//go:build !windows
// +build !windows

package app

import "os/exec"

func PrepareBackgroundCommand(cmd *exec.Cmd) {
}
