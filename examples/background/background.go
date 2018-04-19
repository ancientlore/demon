package main

import (
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("sleep", "5")
	err := cmd.Start() // HL
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait() // HL
	log.Printf("Command finished with error: %v", err)
}
