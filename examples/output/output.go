package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	result, err := exec.Command("bash", "-c", "ls").Output() // HL
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(result))
}
