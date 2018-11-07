package main

import (
	"encoding/json"
	"flag"
)

// Input holds all the user provided flags
type Input struct {
	Reload bool `json:"reload"`
}

// NewInput returns an Input instance containing the user provided flags
func NewInput() *Input {
	reload := flag.Bool("reload", false, "Reload the settings without downtime.")
	flag.Parse()

	return &Input{Reload: *reload}
}

// String returns the json encoded input parameters
func (i Input) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}
