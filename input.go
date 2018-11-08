package main

import (
	"encoding/json"
	"flag"
)

// Input holds all the user provided flags
type Input struct {
	Reload     bool `json:"reload"`
	ShouldQuit bool `json:"should_quit"`
}

// NewInput returns an Input instance containing the user provided flags
func NewInput() *Input {
	input := &Input{}
	input.ShouldQuit = false
	input.Reload = *flag.Bool("reload", false, "Reload the settings without downtime.")

	if input.Reload {
		input.ShouldQuit = true
	}

	flag.Parse()

	return input
}

// String returns the json encoded input parameters
func (i Input) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}
