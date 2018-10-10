package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

// Octaaf is the global bot endpoint
var Octaaf *tgbotapi.BotAPI

func initBot() {
	// Explicitly create this err var, or else Octaaf will be shadowed
	var err error
	Octaaf, err = tgbotapi.NewBotAPI(settings.Telegram.ApiKey)

	if err != nil {
		log.Fatal("Telegram connection error: ", err)
	}

	Octaaf.Debug = settings.Environment == "development"

	log.Info("Authorized on account ", Octaaf.Self.UserName)

	if settings.Environment == "production" {
		sendGlobal(fmt.Sprintf("I'm up and running! 👌\nRunning with version: %v", settings.Version))
		sendGlobal(fmt.Sprintf("Check out the changelog over here: \n%v/tags/%v", GitUri, settings.Version))

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
}

func handle(m *tgbotapi.Message) {
	message := &OctaafMessage{
		m,
		Tracer.StartSpan("Message Received")}

	defer message.Span.Finish()
	message.Span.SetTag("telegram-group-id", message.Chat.ID)
	message.Span.SetTag("telegram-group-name", message.Chat.UserName)
	message.Span.SetTag("telegram-message-id", message.MessageID)
	message.Span.SetTag("telegram-from-id", message.From.ID)
	message.Span.SetTag("telegram-from-username", message.From.UserName)

	go kaliHandler(message)

	isReporter := (message.From.ID == settings.Telegram.ReporterID)
	allowed := true
	if isReporter {
		allowed = message.MessageID%2 == 0
	}

	if message.IsCommand() && !allowed {
		message.Reply("You have been rate limited for spreading lies.")
	} else if message.IsCommand() {
		message.Span.SetOperationName(fmt.Sprintf("Command /%v", message.Command()))
		message.Span.SetTag("telegram-command", message.Command())
		message.Span.SetTag("telegram-command-arguments", message.CommandArguments())
		log.Debugf("Command received: %v", message.Command())
		switch message.Command() {
		case "all":
			all(message)
		case "roll":
			sendRoll(message)
		case "m8ball":
			m8Ball(message)
		case "bodegem":
			sendBodegem(message)
		case "changelog":
			changelog(message)
		case "img", "img_sfw", "more":
			sendImage(message)
		case "stallman":
			sendStallman(message)
		case "search", "search_nsfw":
			search(message)
		case "where":
			where(message)
		case "count":
			count(message)
		case "weather":
			weather(message)
		case "what":
			what(message)
		case "xkcd":
			xkcd(message)
		case "quote":
			quote(message)
		case "next_launch":
			nextLaunch(message)
		case "doubt":
			doubt(message)
		case "issues":
			issues(message)
		case "kalirank":
			kaliRank(message)
		case "iasip":
			iasip(message)
		case "reported":
			reported(message)
		case "remind_me":
			remind(message)
		case "care":
			care(message)
		case "pollentiek":
			pollentiek(message)
		}

		if message.From.ID == settings.Telegram.ModeratorID {
			switch message.Command() {
			case "limit":
				message.Reply("Not implemented yet.")
			}
		}
	}
	if message.MessageID%100000 == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("💯💯💯💯 YOU HAVE MESSAGE %v 💯💯💯💯", message.MessageID))
		msg.ReplyToMessageID = message.MessageID
		msg.ParseMode = "markdown"
		Octaaf.Send(msg)
	}

	// Maintain an array of chat members per group in Redis
	span := message.Span.Tracer().StartSpan(
		"redis /all array",
		opentracing.ChildOf(message.Span.Context()),
	)
	// span := Tracer.StartSpan("redis /all array", message.Span.Context())
	Redis.SAdd(fmt.Sprintf("members_%v", message.Chat.ID), message.From.ID)
	span.Finish()
}

func sendGlobal(message string) {
	msg := tgbotapi.NewMessage(settings.Telegram.KaliID, message)
	_, err := Octaaf.Send(msg)

	if err != nil {
		log.Errorf("Error while sending global '%v': %v", message, err)
	}
}

func reply(message *OctaafMessage, r interface{}) error {
	switch resp := r.(type) {
	default:
		return errors.New("Unkown response type given")
	case string:
		msg := tgbotapi.NewMessage(message.Chat.ID, resp)
		msg.ReplyToMessageID = message.MessageID
		msg.ParseMode = "markdown"
		_, err := Octaaf.Send(msg)
		return err
	case []byte:
		bytes := tgbotapi.FileBytes{Name: "image.jpg", Bytes: resp}
		msg := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)
		msg.Caption = message.CommandArguments()
		msg.ReplyToMessageID = message.MessageID
		_, err := Octaaf.Send(msg)
		return err
	}
}

func getUsername(userID int, chatID int64) (tgbotapi.ChatMember, error) {
	config := tgbotapi.ChatConfigWithUser{
		ChatID:             chatID,
		SuperGroupUsername: "",
		UserID:             userID}

	return Octaaf.GetChatMember(config)
}
