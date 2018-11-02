package main

import (
	"io/ioutil"
	"octaaf/trump"

	"github.com/fogleman/gg"
)

// Assets is an in memory store for the assets folder
type Assets struct {
	Doubt []byte
	Trump *gg.Context
}

func (a *Assets) Load() error {
	var err error
	a.Doubt, err = ioutil.ReadFile("assets/doubt.jpg")

	if err != nil {
		return err
	}

	a.Trump, err = trump.LoadOrder("assets/trump.png", settings.Trump.FontPath, settings.Trump.FontSize)

	return err
}
