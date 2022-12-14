package main

import (
	"fmt"
	"os/exec"
)

func main1() {
	f, err := exec.LookPath("adb")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(f) //  /bin/ls
}
