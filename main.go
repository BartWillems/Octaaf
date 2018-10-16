package main

import (
	"io/ioutil"

	"octaaf/web"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

var settings Settings
var assets Assets

const GitUri = "https://gitlab.com/BartWillems/octaaf"

func main() {

	if _, err := settings.Load(); err != nil {
		log.Fatal("Unable to load/parse the settings file: ", err)
	}

	settings.Version = loadVersion()

	if settings.Environment != "production" {
		log.SetLevel(log.DebugLevel)
	}

	err := assets.Load()
	if err != nil {
		log.Fatalf("Unable to load the assets: %v", err)
	}

	initRedis()

	err = initDB()
	if err != nil {
		log.Fatalf("Couldn't establish a database connection: %v", err)
	}

	err = migrateDB()
	if err != nil {
		log.Fatalf("DB Migration error: %v", err)
	}

	err = initBot()
	if err != nil {
		log.Fatalf("Telegram connection error: %v", err)
	}

	go loadReminders()

	cron := initCrons()
	cron.Start()
	defer cron.Stop()

	go func() {
		router := web.New(web.Connections{
			Octaaf:      Octaaf,
			Postgres:    DB,
			Redis:       Redis,
			KaliID:      settings.Telegram.KaliID,
			Environment: settings.Environment,
		})
		err := router.Run()
		if err != nil {
			log.Errorf("Gin creation error: %v", err)
		}
	}()

	closer := initJaeger("octaaf")
	defer closer.Close()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := Octaaf.GetUpdatesChan(u)

	if err != nil {
		log.Fatalf("Failed to fetch updates: %v", err)
	}

	for update := range updates {

		if update.Message == nil {
			continue
		}

		go handle(update.Message)
	}
}

func loadVersion() string {
	bytes, err := ioutil.ReadFile("assets/version")

	if err != nil {
		log.Errorf("Error while loading version string: %v", err)
		return ""
	}
	log.Infof("Loaded version %v", string(bytes))
	return string(bytes)
}
