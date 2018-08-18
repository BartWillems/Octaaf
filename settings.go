package main

import "github.com/BurntSushi/toml"

type Settings struct {
	Environment string
	Version     string

	Telegram telegram
	Database database
	Redis    redis
	Google   google
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

func (c *Settings) Load() (toml.MetaData, error) {
	return toml.DecodeFile("config/settings.toml", c)
}