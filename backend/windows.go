//go:build windows
// +build windows

package app

import (
	"os/exec"
	"syscall"
)

func PrepareBackgroundCommand(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

// 同步执行命令,但是是使用的cmd /c的形式
func RunCommandWithEnv(commands string) (stdout string, stderr string, exitCode int) {
	log.Debug("run command:", commands)

	var cmd *exec.Cmd

	cmd = exec.Command("cmd.exe")
	//核心点,直接修改执行命令方式
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`/c %s`, commands)}

	PrepareBackgroundCommand(cmd)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()

	stdout = outbuf.String()
	stderr = errbuf.String()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			log.Infof("Could not get exit code for failed program: %v", commands)
			exitCode = 0
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	// log.Printf("command result, stdout: %v, stderr: %v, exitCode: %v", stdout, stderr, exitCode)

	// args := strings.Split(command, " ")

	return
}
