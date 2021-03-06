package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		panic("args are not provided: <env_dir> <command ...>")
	}
	dir, command := os.Args[1], os.Args[2:]
	envs, err := ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	exitCode := RunCmd(command, envs)
	os.Exit(exitCode)
}
