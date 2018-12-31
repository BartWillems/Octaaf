package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Env(t *testing.T) {
	var settings Settings

	os.Setenv("TELEGRAM_API_KEY", "abc")
	os.Setenv("KALI_ID", "123")
	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("REDIS_DB", "5")

	err := settings.Load()

	assert.Nil(t, err)
	assert.Equal(t, settings.Environment, "production")
	assert.Equal(t, settings.Telegram.APIKEY, "abc")
	assert.Equal(t, settings.Telegram.KaliID, int64(123))
	assert.Equal(t, settings.Redis.DB, int(5))
}
