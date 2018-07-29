package main

import (
	"log"

	"github.com/gobuffalo/envy"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	envy.Load("config/.env")

	connectDB()
	migrateDB()
	initRedis()
	initBot()

	initCrons()

	defer Cron.Stop()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := Octaaf.GetUpdatesChan(u)

	if err != nil {
		log.Panicf("Failed to fetch updates: %v", err)
	}

	for update := range updates {

		if update.Message == nil {
			continue
		}

		go handle(update.Message)
	}
}
