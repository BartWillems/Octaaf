package main

import (
	"fmt"
	"math/rand"
	"octaaf/models"
	"time"
	log "github.com/sirupsen/logrus"
)

// KaliCount is an integer that holds the ID of the last sent message in the Kali group
var KaliCount int

func kaliHandler(message *OctaafMessage) {
	if message.Chat.ID == settings.Telegram.KaliID {
		log.Debug("Kalimember found")
		KaliCount = message.MessageID
	}
}

func setKaliCount() {
	if KaliCount <= 0 {
		log.Error("Unable to save today's KaliCount because it's ", KaliCount)
		return
	}
	count := models.MessageCount{
		Count: KaliCount,
		Diff:  0,
	}

	err := DB.Save(&count)
	if err != nil {
		log.Error("Unable to save today's kalicount: ", err)
	}
}

func checkIn() {
	rand.Seed(time.Now().UnixNano())
	// Generate a random time moment
	randomTime, _ := time.ParseDuration(fmt.Sprintf("%vs", rand.Intn(24*3600)))
	time.Sleep(randomTime)
	sendGlobal("RANDOM CHECKIN!!!! T-60s")
	time.Sleep(60 * time.Second)
	sendGlobal("Checking complete.")
}
