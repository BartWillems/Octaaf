package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

// Input holds all the user provided flags
type Input struct {
	Reload     bool `json:"reload"`
	ShouldQuit bool `json:"should_quit"`
}

// NewInput returns an Input instance containing the user provided flags
func NewInput() *Input {
	reload := flag.Bool("reload", false, "Reload the settings without downtime.")
	flag.Parse()

	input := &Input{
		Reload:     *reload,
		ShouldQuit: false,
	}

	if input.Reload {
		input.ShouldQuit = true
	}

	return input
}

// String returns the json encoded input parameters
func (i Input) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

func (i *Input) TriggerReload() error {
	_, err := http.Post(fmt.Sprintf("http://%s/%s/reload", SOCKET_URI, SOCKET_PATH), "", nil)
	return err
}
