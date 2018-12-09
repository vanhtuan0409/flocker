package main

import (
	"fmt"
	"os"
)

func main() {
	cmd := os.Args[1]

	switch cmd {
	case "run":
		fmt.Println("run")
	default:
		fmt.Println("Unknown command")
	}
}
