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

// Start starts a process after processing environment variables in the command line
// and connecting pipes for I/O.
func (p *process) Start(dest, errDest io.Writer) error {
	var err error

	name := p.Command[0]
	args := make([]string, len(p.Command)-1)
	copy(args, p.Command[1:])
	for i := range args {
		if strings.HasPrefix(args[i], "$") {
			args[i] = os.Getenv(strings.TrimPrefix(args[i], "$"))
		}
	}
	log.Print("Starting ", name, " ", strings.Join(args, " "))
	p.cmd = exec.Command(name, args...)
	p.cmd.Dir = p.Dir

	p.stdin, err = p.cmd.StdinPipe()
	if err != nil {
		return err
	}
	p.stdout, err = p.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	p.wg.Add(1)
	go readPipe(&p.wg, p.stdout, dest)
	p.cmd.Stderr = p.cmd.Stdout
	/*
		p.stderr, err = p.cmd.StderrPipe()
		if err != nil {
			return err
		}
		p.wg.Add(1)
		go readPipe(&p.wg, p.stderr, errDest)
	*/

	err = p.cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

// Wait completes the execution of a process and waits for it to finish.
func (p *process) Wait() error {
	defer func() {
		p.stderr = nil
		p.stdout = nil
		p.stdin = nil
		p.cmd = nil
	}()
	if p.stdin != nil {
		if p.ExitInput != "" {
			_, err := p.Write([]byte(p.ExitInput + "\n"))
			if err != nil {
				log.Print(err)
			}
			time.Sleep(100 * time.Millisecond)
		}
		p.stdin.Close()
	}
	if p.cmd != nil {
		log.Print("Waiting for ", p.Command[0])
		p.wg.Wait()
		return p.cmd.Wait()
	}
	return nil
}

// Write is ised to write information to stdin.
func (p *process) Write(b []byte) (int, error) {
	log.Print(p.Title, ": ", string(b))
	return p.stdin.Write(b)
}

// Interrupt sends an interrupt signal to the process.
func (p *process) Interrupt() error {
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
