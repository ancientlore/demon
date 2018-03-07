package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
)

// site represents a web site in the view.
type site struct {
	Title string `json:"title"` // Title to display
	URL   string `json:"url"`   // URL of the site
}

// process represents a OS process that commands can be sent to.
type process struct {
	Title   string   `json:"title"`         // Title to display
	Command []string `json:"command"`       // Command to run to start the process
	Dir     string   `json:"dir,omitempty"` // Working directory

	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	wg     sync.WaitGroup
}

// step is a step of the demo that includes a command to run.
type step struct {
	Title     string `json:"title"`               // Title to display
	Desc      string `json:"desc,omitempty"`      // Description to display (Markdown allowed)
	ID        string `json:"id"`                  // ID of the step to operate on
	Input     string `json:"input"`               // Command to send to process
	Interrupt bool   `json:"interrupt,omitempty"` // Send to interrupt command
}

// config holds the sites, processes, and steps for the demo.
type config struct {
	Sites     map[string]*site    `json:"sites,omitempty"`     // Sites to show
	Processes map[string]*process `json:"processes,omitempty"` // Processes to start
	Steps     []step              `json:"steps"`               // Steps to enable
}

// readConfig loads a JSON file that stores a config.
func readConfig(file string) (*config, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var c config
	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Start starts all of the processes in a configuration.
func (c *config) Start() error {
	var merr multiError
	for _, p := range c.Processes {
		err := p.Start(os.Stdout, os.Stderr)
		if err != nil {
			merr = append(merr, err)
		}
	}
	if len(merr) > 0 {
		c.Wait()
		return merr
	}
	return nil
}

// Wait completes all of the processes in a configuration.
func (c *config) Wait() error {
	var merr multiError
	for _, p := range c.Processes {
		err := p.Wait()
		if err != nil {
			merr = append(merr, err)
		}
	}
	if len(merr) > 0 {
		return merr
	}
	return nil
}
