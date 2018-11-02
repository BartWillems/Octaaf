package main

import (
	"image"
	"io/ioutil"

	"github.com/fogleman/gg"
)

// Assets is an in memory store for the assets folder
type Assets struct {
	Doubt []byte
	Trump image.Image
}

func (a *Assets) Load() error {
	var err error
	a.Doubt, err = ioutil.ReadFile("assets/doubt.jpg")

	if err != nil {
		return err
	}

	a.Trump, err = gg.LoadPNG("assets/trump.png")

	return err
}
