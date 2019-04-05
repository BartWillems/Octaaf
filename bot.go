package main

import (
	"fmt"
	"octaaf/jaeger"
	"octaaf/kcoin"
	"octaaf/kcoin/rewards"
	"octaaf/models"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

// Octaaf is the global bot endpoint
var Octaaf *tgbotapi.BotAPI

func initBot() error {
	// Explicitly create this err var, or else Octaaf will be shadowed
	var err error
	Octaaf, err = tgbotapi.NewBotAPI(settings.Telegram.APIKEY)

	if err != nil {
		return err
	}

	Octaaf.Debug = settings.Environment == development

	log.Info("Authorized on account ", Octaaf.Self.UserName)

	if settings.Environment == "production" {
		sendGlobal(fmt.Sprintf("I'm up and running! 👌\nRunning with version: %v", Version))
		sendGlobal(fmt.Sprintf("Check out the changelog over here: \n%v/tags/%v", GitURI, Version))

		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			sendGlobal("I'm going to sleep! 💤💤")
			DB.Close()
			Redis.Close()
			os.Exit(0)
		}()
	}

	return nil
}

func handle(m *tgbotapi.Message) {
	message := &OctaafMessage{
		m,
		jaeger.Tracer.StartSpan("Message Received"),
		true,  // IsMarkdown
		false, // KeyboardCloser
	}

	defer message.Span.Finish()
	message.Span.SetTag("telegram-group-id", message.Chat.ID)
	message.Span.SetTag("telegram-group-name", message.Chat.Title)
	message.Span.SetTag("telegram-message-id", message.MessageID)
	message.Span.SetTag("telegram-from-id", message.From.ID)
	message.Span.SetTag("telegram-from-username", message.From.UserName)

	go kaliHandler(message)

	var err error
	transactionSucceeded := true
	if settings.Kalicoin.Enabled {
		transactionSucceeded, err = kcoin.HandleTransaction(message.Message, message.Span)

		if time.Now().Hour() == 13 && time.Now().Minute() == 37 && message.Text == "1337" {
			go rewards.StoreUser(Redis, message.Chat.ID, message.From.ID, "kalivent")
		}

		if time.Now().Hour() == 16 && time.Now().Minute() == 20 && message.Text == "420" {
			go rewards.StoreUser(Redis, message.Chat.ID, message.From.ID, "kalivent")
		}
	}

	if !transactionSucceeded {
		message.Reply(err.Error())
	} else if message.IsCommand() {
		executeCommand(message)
	}

	if message.MessageID%1e5 == 0 {
		message.Reply(fmt.Sprintf("💯💯💯💯 YOU HAVE MESSAGE %v 💯💯💯💯", message.MessageID))
	}

	// Maintain an array of chat members per group in Redis
	span := message.Span.Tracer().StartSpan(
		"redis /all array",
		opentracing.ChildOf(message.Span.Context()),
	)
	Redis.SAdd(fmt.Sprintf("members_%v", message.Chat.ID), message.From.ID)
	span.Finish()
}

func executeCommand(message *OctaafMessage) error {
	message.Span.SetOperationName(fmt.Sprintf("Command /%v", message.Command()))
	message.Span.SetTag("is-command", true)
	message.Span.SetTag("telegram-command", message.Command())
	message.Span.SetTag("telegram-command-arguments", message.CommandArguments())
	switch message.Command() {
	case "all":
		return all(message)
	case "roll":
		return sendRoll(message)
	case "m8ball":
		return m8Ball(message)
	case "bodegem":
		return sendBodegem(message)
	case "changelog":
		return changelog(message)
	case "img", "img_sfw", "more", "img_censored":
		return sendImage(message)
	case models.ImgQuote, models.VodQuote, models.AudioQuote:
		return msgQuote(message)
	case "stallman":
		return sendStallman(message)
	case "search", "search_nsfw":
		return search(message)
	case "where":
		return where(message)
	case "count":
		return count(message)
	case "weather":
		return weather(message)
	case "what":
		return what(message)
	case "xkcd":
		return xkcd(message)
	case "quote", "presidential_quote":
		return quote(message)
	case "next_launch":
		return nextLaunch(message)
	case "issues":
		return gitlabIssues(message)
	case "doubt":
		return doubt(message)
	case "kalirank":
		return kaliRank(message)
	case "remind_me":
		return remind(message)
	case "care":
		return care(message)
	case "pollentiek":
		return pollentiek(message)
	case "presidential_order":
		return presidentialOrder(message)
	case "wallet":
		if !settings.Kalicoin.Enabled {
			return message.Reply("Kalicoin has been disabled by the administrator")
		}
		wallet, err := kcoin.Wallet(message.Message, message.Span)

		if err != nil {
			return message.Reply(err.Error())
		}

		return message.Reply(fmt.Sprintf("💰Balance: %v", wallet.Capital))
	}

	return nil
}

func sendGlobal(message string) {
	msg := tgbotapi.NewMessage(settings.Telegram.KaliID, message)
	msg.ParseMode = "markdown"
	_, err := Octaaf.Send(msg)

	if err != nil {
		log.Errorf("Error while sending global '%v': %v", message, err)
	}
}

func getUser(userID int, chatID int64) (tgbotapi.ChatMember, error) {
	config := tgbotapi.ChatConfigWithUser{
		ChatID:             chatID,
		SuperGroupUsername: "",
		UserID:             userID}

	return Octaaf.GetChatMember(config)
}

// Returns a username for a user ID and a chat ID
func getUserName(userID int, chatID int64) (string, error) {
	user, err := getUser(userID, chatID)

	if err != nil {
		return "", err
	}

	return user.User.UserName, nil
}

func getUserNames(userIDS []int, chatID int64) []string {
	var users []string

	var wg sync.WaitGroup
	// Get the members' usernames
	for _, userID := range userIDS {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			username, _ := getUserName(userID, chatID)
			users = append(users, username)
		}(userID)
	}

	wg.Wait()
	return users
}
