package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := os.Args[1]

	switch cmd {
	case "run":
		run(os.Args[2], os.Args[3:])
	default:
		fmt.Println("Unknown command")
	}
}

func run(command string, opts []string) {
	cmd := exec.Command(command, opts...)
	bindCmdStdio(cmd)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}
	must(cmd.Run())
}

func bindCmdStdio(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
