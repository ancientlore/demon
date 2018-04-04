package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// site represents a web site in the view.
type site struct {
	Title    string `json:"title"`    // Title to display
	URL      string `json:"url"`      // URL of the site
	Position int    `json:"position"` // Position of the panel
}

// process represents a OS process that commands can be sent to.
type process struct {
	Title     string   `json:"title"`         // Title to display
	Command   []string `json:"command"`       // Command to run to start the process
	Dir       string   `json:"dir,omitempty"` // Working directory
	Position  int      `json:"position"`      // Position of the panel
	ExitInput string   `json:"exitInput"`     // Text to send at exit

	cmd    *exec.Cmd      // Command object for running the process
	stdin  io.WriteCloser // Input stream for the process
	stdout io.ReadCloser  // Output stream for the process
	stderr io.ReadCloser  // Error stream for the process
	wg     sync.WaitGroup // Wait group used to make sure the readers are done
	wm     sync.Mutex     // Mutex to make sure writes to stdin are serialized
}

// step is a step of the demo that includes a command to run.
type step struct {
	Title     string `json:"title"`               // Title to display
	Desc      string `json:"desc,omitempty"`      // Description to display (Markdown allowed)
	ID        string `json:"id"`                  // ID of the process to operate on
	Input     string `json:"input"`               // Command to send to process
	Interrupt bool   `json:"interrupt,omitempty"` // Send to interrupt command
}

// config holds the sites, processes, and steps for the demo.
type config struct {
	Title         string              `json:"title"`               // Title to display
	StepsPosition int                 `json:"stepsPosition"`       // Position of the steps panel
	Sites         map[string]*site    `json:"sites,omitempty"`     // Sites to show
	Processes     map[string]*process `json:"processes,omitempty"` // Processes to start
	Steps         []step              `json:"steps"`               // Steps to enable

	Loop []int `json:"-"` // Used to iterate in position order in the template
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
	err = c.Validate()
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

// Validate checks for inconsistencies in the config.
func (c *config) Validate() error {
	var m multiError
	for i := range c.Steps {
		if _, ok := c.Processes[c.Steps[i].ID]; !ok && c.Steps[i].Input != "" {
			m = append(m, fmt.Errorf("Step %d %q references missing process %q", i, c.Steps[i].Title, c.Steps[i].ID))
		}
	}
	count := len(c.Processes) + len(c.Sites) + 1
	c.Loop = make([]int, count)
	for i := range c.Loop {
		c.Loop[i] = i
	}
	for k, v := range c.Processes {
		if v.Position < 0 || v.Position >= count {
			m = append(m, fmt.Errorf("Position of process %q is out of bounds: %d", k, v.Position))
		}
	}
	for k, v := range c.Sites {
		if v.Position < 0 || v.Position >= count {
			m = append(m, fmt.Errorf("Position of site %q is out of bounds: %d", k, v.Position))
		}
		if strings.HasPrefix(v.URL, "$") {
			v.URL = os.Getenv(strings.TrimPrefix(v.URL, "$"))
		}
	}
	if c.StepsPosition < 0 || c.StepsPosition >= count {
		m = append(m, fmt.Errorf("Position of steps panel is out of bounds: %d", c.StepsPosition))
	}
	if len(m) > 0 {
		return m
	}

	return nil
}
