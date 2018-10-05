package main

import "io/ioutil"

// Assets is an in memory store for the assets folder
type Assets struct {
	Doubt []byte
}

func (a *Assets) Load() error {
	var err error
	a.Doubt, err = ioutil.ReadFile("assets/doubt.jpg")

	return err
}
