package main

import (
	"octaaf/trump"

	"github.com/BurntSushi/toml"
)

type Settings struct {
	Environment string

	Telegram telegram
	Database database
	Redis    redis
	Google   google
	Jaeger   jaeger
	Trump    trump.TrumpConfig
}

type telegram struct {
	ApiKey     string `toml:"api_key"`
	KaliID     int64  `toml:"kali_id"`
	ReporterID int    `toml:"reporter_id"`
}

type database struct {
	Uri string `toml:"uri"`
}

type redis struct {
	Uri      string `toml:"uri"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

type google struct {
	ApiKey string `toml:"api_key"`
}

type jaeger struct {
	ServiceName string `toml:"service_name"`
}

// Load parses the toml file and returns a Settings struct
func (c *Settings) Load() (toml.MetaData, error) {
	return toml.DecodeFile("config/settings.toml", c)
}
