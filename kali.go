package main

import (
	"fmt"
	"math/rand"
	"octaaf/models"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// KaliCount is an integer that holds the ID of the last sent message in the Kali group
var KaliCount int

// IsCheckinTime checks if the random checkin is enabled
var IsCheckinTime = false

// KaliCheckersKey is the redis key for the kali random checkin members
const KaliCheckersKey = "kalicheckers"

func kaliHandler(message *OctaafMessage) {
	if message.Chat.ID != settings.Telegram.KaliID {
		return
	}

	log.Debug("Kalimember found")
	KaliCount = message.MessageID

	if IsCheckinTime && (strings.Contains(message.Text, "check") || strings.ContainsAny(message.Text, "✅✔️☑️")) {
		log.Infof("Random checker found: %v", message.From.ID)
		Redis.SAdd(KaliCheckersKey, message.From.ID)
	}

	if message.IsCommand() {
		if message.Command() == "checkrepublic" {
			getKaliCheckers(message)
		}
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

func saveKaliCheckers() {
	checkers := Redis.SMembers(KaliCheckersKey).Val()

	for _, checker := range checkers {
		UserID, err := strconv.Atoi(checker)

		if err != nil {
			log.Errorf("Unable to convert kalickecker's user ID to int: %v", err)
			continue
		}

		err = DB.Save(&models.Kalichecker{
			UserID: UserID,
		})

		if err != nil {
			log.Errorf("Unable to store kalichecker: %v", err)
		}
	}

	Redis.Del(KaliCheckersKey)
}

func getKaliCheckers(message *OctaafMessage) error {
	var KalicheckerStats models.KalicheckerStats
	err := KalicheckerStats.Top(DB)

	if err != nil || len(KalicheckerStats) == 0 {
		log.Errorf("KALI ERROR: %v", err)
		return message.Reply("404 - No entries found. Maybe tomorrow?")
	}

	response := "*Rank: count - name*\n"

	for index, stat := range KalicheckerStats {
		username, err := getUserName(stat.UserID, settings.Telegram.KaliID)

		if err != nil {
			log.Errorf("Unable to fetch username for the kalicheckerstats: %v", err)
			message.Span.SetTag("error", err)
			continue
		}

		response += fmt.Sprintf("*%v:* %v - @%v \n", index+1, stat.Count, username)
	}

	return message.Reply(response)

}

func checkIn() {
	rand.Seed(time.Now().UnixNano())
	// Generate a random time moment
	randomTime, _ := time.ParseDuration(fmt.Sprintf("%vs", rand.Intn(24*3600)))
	time.Sleep(randomTime)

	sendGlobal("RANDOM CHECKIN!!!! T-60s")
	IsCheckinTime = true

	time.Sleep(60 * time.Second)

	sendGlobal("Checking complete.")
	IsCheckinTime = false
	saveKaliCheckers()
}
