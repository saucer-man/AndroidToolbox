package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// 带有outStr和errStr和exitCode的命令执行
func RunCommand(name string, args ...string) (stdout string, stderr string, exitCode int) {
	log.Debug("run command:", name, args)
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

// 同步执行命令,但是是使用的cmd /c的形式
func RunCommandWithEnv(commands string) (stdout string, stderr string, exitCode int) {
	log.Debug("run command:", commands)

	sysType := runtime.GOOS
	var cmd *exec.Cmd
	if sysType == "windows" {
		cmd = exec.Command("cmd.exe")
		//核心点,直接修改执行命令方式
		cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`/c %s`, commands), HideWindow: true}
	} else {
		cmd = exec.Command("bash")
		//核心点,直接修改执行命令方式
		cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`-c %s`, commands), HideWindow: true}
	}

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

// 同步执行命令
func Excute(commands []string) {
	log.Infof("将要执行%+v", commands)

	command := commands[0]
	args := []string{}
	if len(command) > 1 {
		args = commands[1:]
	}
	// args := strings.Split(command, " ")

	stdout, stderr, exitCode := RunCommand(command, args...)
	log.Infof("执行结果为,stdout:%+v, stderr:%+v, exitCode:%+v", stdout, stderr, exitCode)
}

// 同步执行命令
func ExcuteWithEnv(commands string) {
	log.Infof("将要执行%+v", commands)
	// args := strings.Split(command, " ")

	stdout, stderr, exitCode := RunCommandWithEnv(commands)
	log.Infof("执行结果为,stdout:%+v, stderr:%+v, exitCode:%+v", stdout, stderr, exitCode)
}

func main1() {
	filePath := "C:\\Users\\yanq\\Downloads\\CheatEngine75.exe"
	toSavePath := "/storage/emulated/0/MT2/"
	ExcuteWithEnv(fmt.Sprintf("adb -s 1B191FDF60056H push  \"%s\" \"%s\"", filePath, toSavePath))
}
