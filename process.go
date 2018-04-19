package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
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
func (p *process) Start(dest, errDest io.Writer) error {
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
	p.wg.Add(1)
	go readPipe(&p.wg, p.stdout, dest) // HL
	p.cmd.Stderr = p.cmd.Stdout
	/*
		p.stderr, err = p.cmd.StderrPipe()
		if err != nil {
			return err
		}
		p.wg.Add(1)
		go readPipe(&p.wg, p.stderr, errDest)
	*/

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
		p.stderr = nil
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
		p.wg.Wait()
		return p.cmd.Wait() // HL
	}
	return nil
}

// Write is ised to write information to stdin.
func (p *process) Write(b []byte) (int, error) {
	p.wm.Lock()
	defer p.wm.Unlock()
	log.Print(p.Title, ": Write: ", string(b))
	return p.stdin.Write(b)
}

// Interrupt sends an interrupt signal to the process.
func (p *process) Interrupt() error {
	log.Print(p.Title, ": Interrupt")
	return p.cmd.Process.Signal(os.Interrupt)
}

// readPipe is used as a gorotuine to process data coming back from the process.
func readPipe(wg *sync.WaitGroup, p io.Reader, dest io.Writer) {
	defer wg.Done()

	scanner := bufio.NewScanner(p)
	for scanner.Scan() {
		fmt.Fprintln(dest, scanner.Text())
		// fmt.Fprintln(os.Stdout, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Print(err)
	}
}
