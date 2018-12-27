package main

import (
	"octaaf/web"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

var settings Settings
var assets Assets

// Version is a git tag that get's added at compile time
var Version string

const GitUri = "https://gitlab.com/BartWillems/octaaf"

func main() {

	if _, err := settings.Load(); err != nil {
		log.Fatal("Unable to load/parse the settings file: ", err)
	}

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
		log.Info("Starting the external api...")
		router := web.New(web.Connections{
			Octaaf:      Octaaf,
			Postgres:    DB,
			Redis:       Redis,
			KaliID:      settings.Telegram.KaliID,
			Environment: settings.Environment,
			TrumpCfg:    &settings.Trump,
			Trump:       &assets.Trump,
		})
		err := router.Run()
		if err != nil {
			log.Errorf("External API creation error: %v", err)
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
