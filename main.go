package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
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

	volumeDir := genContainerDir()
	fmt.Println(volumeDir)
	must(syscall.Chroot(volumeDir))
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	must(cmd.Run())
}

func genContainerDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Panicf("Cannot get current dir. ERR: %v\n", err)
	}

	containerDir := path.Join(pwd, "container")

	volumeDir := path.Join(containerDir, randString(10))
	for {
		_, err := os.Stat(volumeDir)
		if err != nil && os.IsNotExist(err) {
			break
		}
		volumeDir = path.Join(pwd, "container", randString(10))
	}

	err = os.MkdirAll(volumeDir, os.ModePerm)
	if err != nil {
		log.Panicf("Cannot create volume dir. ERR: %v\n", err)
	}

	return volumeDir
}
