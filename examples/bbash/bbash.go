package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("bash")

	stdin, err := cmd.StdinPipe() // HL
	if err != nil {
		log.Fatal(err)
	}

	stdout, err := cmd.StdoutPipe() // HL
	if err != nil {
		log.Fatal(err)
	}

	// Combine stderr and stdout
	cmd.Stderr = cmd.Stdout // HL

	err = cmd.Start() // HL
	if err != nil {
		log.Fatal(err)
	}

	// launch routine to write commands
	go func() {
		defer stdin.Close() // HL
		io.WriteString(stdin, "echo stdout; echo 1>&2 stderr; echo A LIST OF FILES\n")
		io.WriteString(stdin, "ls\n")
		io.WriteString(stdin, "exit\n")
	}()

	// Read the output
	s := bufio.NewScanner(stdout)
	for s.Scan() {
		fmt.Println(s.Text())
	}
	if err = s.Err(); err != nil {
		log.Fatal(err)
	}

	// Wait on command completion (after output read)
	err = cmd.Wait() // HL
	if err != nil {
		log.Fatal(err)
	}
}
