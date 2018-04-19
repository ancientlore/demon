package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("cat")
	stdin, err := cmd.StdinPipe() // HL
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close() // HL
		io.WriteString(stdin, "values written to stdin are passed to cmd's standard input")
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", out)
}
