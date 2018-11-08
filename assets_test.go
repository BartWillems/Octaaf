package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Assets(t *testing.T) {
	var assets Assets
	var empty []byte
	err := assets.Load()
	assert.Nil(t, err)
	assert.NotEqual(t, assets.Doubt, empty)
}
