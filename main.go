package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

var (
	internal = flag.Bool("internal", false, "Mark process as internal isolation")
)

func main() {
	flag.Parse()

	cmd := os.Args[1]
	switch cmd {
	case "run":
		setupIsolation(os.Args[2], os.Args[3:])
	case "internal":
		run(os.Args[2], os.Args[3:])
	default:
		fmt.Println("Unknown command")
	}
}

func setupIsolation(command string, opts []string) {
	appendedArgs := append([]string{"internal"}, os.Args[2:]...)
	cmd := exec.Command("/proc/self/exe", appendedArgs...)
	bindCmdStdio(cmd)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}
	must(cmd.Run())
}

func run(command string, opts []string) {
	fmt.Printf("Running command %s as PID %d\n", command, os.Getpid())
	cmd := exec.Command(command, opts...)
	bindCmdStdio(cmd)
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
