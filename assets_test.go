package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Assets(t *testing.T) {
	var assets Assets
	err := assets.Load()
	assert.Nil(t, err)
}
