package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Clone makes a copy of the process, ready to be started.
func (p *process) Clone() *process {
	return &process{
		Title:     p.Title,
		Command:   p.Command,
		Dir:       p.Dir,
		Position:  p.Position,
		ExitInput: p.ExitInput,
	}
}

// Start starts a process after processing environment variables in the command line
// and connecting pipes for I/O.
func (p *process) Start() error {
	var err error

	if p.cmd != nil {
		return fmt.Errorf(p.Title, ": Process already started")
	}

	name := p.Command[0]
	args := make([]string, len(p.Command)-1)
	copy(args, p.Command[1:])
	for i := range args {
		if strings.HasPrefix(args[i], "$") {
			args[i] = os.Getenv(strings.TrimPrefix(args[i], "$"))
		}
	}
	log.Print(p.Title, ": Starting: ", name, " ", strings.Join(args, " "))
	p.cmd = exec.Command(name, args...) // HL
	p.cmd.Dir = p.Dir

	p.stdin, err = p.cmd.StdinPipe() // HL
	if err != nil {
		return err
	}
	p.stdout, err = p.cmd.StdoutPipe() // HL
	if err != nil {
		return err
	}
	p.cmd.Stderr = p.cmd.Stdout // HL

	err = p.cmd.Start() // HL
	if err != nil {
		return err
	}

	return nil
}

// Wait completes the execution of a process and waits for it to finish.
func (p *process) Wait() error {
	p.wm.Lock()
	defer p.wm.Unlock()
	defer func() {
		p.stdout = nil
		p.stdin = nil
		p.cmd = nil
	}()
	if p.stdin != nil {
		if p.ExitInput != "" {
			log.Print(p.Title, ": Writing exit: ", p.ExitInput)
			_, err := p.stdin.Write([]byte(p.ExitInput + "\n"))
			if err != nil {
				log.Print(p.Title, ": Error writing exit: ", err)
			}
			time.Sleep(100 * time.Millisecond)
		}
		log.Print(p.Title, ": Closing stdin")
		p.stdin.Close() // HL
	}
	if p.cmd != nil {
		log.Print(p.Title, ": Starting wait: ", p.Command[0])
		return p.cmd.Wait() // HL
	}
	return nil
}

// Write is used to write information to stdin.
func (p *process) Write(b []byte) (int, error) {
	p.wm.Lock()
	defer p.wm.Unlock()
	log.Print(p.Title, ": Write: ", string(b))
	return p.stdin.Write(b)
}

// Read is used to read information from stdout.
func (p *process) Read(b []byte) (n int, err error) {
	return p.stdout.Read(b)
}

// Interrupt sends an interrupt signal to the process.
func (p *process) Interrupt() error {
	log.Print(p.Title, ": Interrupt")
	return p.cmd.Process.Signal(os.Interrupt)
}
