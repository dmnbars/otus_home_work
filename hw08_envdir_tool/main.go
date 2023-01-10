package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Printf("not enough arguments")
		os.Exit(1)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		os.Exit(1)
	}

	os.Exit(RunCmd(os.Args[2:], env))
}
