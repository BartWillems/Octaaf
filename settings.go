package main

import (
	"encoding/json"
	"octaaf/jaeger"
	"octaaf/kcoin"
	"octaaf/trump"

	env "github.com/BartWillems/go-env"
	"github.com/BurntSushi/toml"
	"github.com/imdario/mergo"
	log "github.com/sirupsen/logrus"
)

// Settings is the struct that holds the runtime config
type Settings struct {
	Environment string `env:"ENVIRONMENT"`

	Kalicoin kcoin.Config
	Telegram telegram
	Database database
	Redis    redis
	Google   google
	Jaeger   jaeger.Config
	Trump    trump.Config
}

type telegram struct {
	APIKEY string `toml:"api_key" env:"TELEGRAM_API_KEY"`
	KaliID int64  `toml:"kali_id" env:"KALI_ID"`
}

type database struct {
	URI string `toml:"uri" env:"DATABASE_URI"`
}

type redis struct {
	URI      string `toml:"uri" env:"REDIS_URI"`
	Password string `toml:"password" env:"REDIS_PASSWORD"`
	DB       int    `toml:"db" env:"REDIS_DB"`
}

type google struct {
	APIKEY string `toml:"api_key" env:"GOOGLE_API_KEY"`
}

// Load parses the toml file and returns a Settings struct
func (s *Settings) Load() error {
	_, err := toml.DecodeFile("config/settings.toml", s)
	if err != nil {
		log.Warnf("Unable to load settings file: %v, falling back to environment variables.", err)
	}

	var envSettings Settings
	_, err = env.UnmarshalFromEnviron(&envSettings)

	if err != nil {
		return err
	}

	// Merge the settings loaded from the envrionment into the config file settings
	if err := mergo.Merge(s, envSettings, mergo.WithOverride); err != nil {
		return err
	}

	return nil
}

// String returns the settings struct as a json string
func (s Settings) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}
