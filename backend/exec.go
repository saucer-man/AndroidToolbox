package app

import (
	"bytes"
	"syscall"

	"os/exec"

	log "github.com/sirupsen/logrus"
)

// 测试某个命令是否存在
func CheckExists(binName string) {
	path, err := exec.LookPath(binName)
	if err != nil {
		log.Fatal("didn't find 'adb' executable\n")
	} else {
		log.Info("'adb' executable is in '%s'\n", path)
	}
}

// 带有outStr和errStr的命令执行
func Exex(command string, args ...string) (string, string) {
	cmd := exec.Command(command, args...)
	var stdin, stdout, stderr bytes.Buffer
	cmd.Stdin = &stdin
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	outStr, errStr := stdout.String(), stderr.String()
	log.Info("out:\n%s\nerr:\n%s\n", outStr, errStr)
	return outStr, errStr
}

// 带有outStr和errStr和exitCode的命令执行
func RunCommand(name string, args ...string) (stdout string, stderr string, exitCode int) {
	log.Info("run command:", name, args)
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(name, args...)
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
			log.Infof("Could not get exit code for failed program: %v, %v", name, args)
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
	return
}
