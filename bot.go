package main

import (
	"fmt"
	"octaaf/markdown"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

// Octaaf is the global bot endpoint
var Octaaf *tgbotapi.BotAPI

func initBot() error {
	// Explicitly create this err var, or else Octaaf will be shadowed
	var err error
	Octaaf, err = tgbotapi.NewBotAPI(settings.Telegram.ApiKey)

	if err != nil {
		return err
	}

	Octaaf.Debug = settings.Environment == "development"

	log.Info("Authorized on account ", Octaaf.Self.UserName)

	if settings.Environment == "production" {
		sendGlobal(fmt.Sprintf("I'm up and running! 👌\nRunning with version: %v", Version))
		sendGlobal(fmt.Sprintf("Check out the changelog over here: \n%v/tags/%v", GitUri, Version))

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
		Tracer.StartSpan("Message Received"),
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

	if message.IsCommand() {
		message.Span.SetOperationName(fmt.Sprintf("Command /%v", message.Command()))
		message.Span.SetTag("is-command", true)
		message.Span.SetTag("telegram-command", message.Command())
		message.Span.SetTag("telegram-command-arguments", message.CommandArguments())
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
		case "img", "img_sfw", "more", "img_censored":
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
		case "quote", "presidential_quote":
			quote(message)
		case "next_launch":
			nextLaunch(message)
		case "doubt":
			doubt(message)
		case "kalirank":
			kaliRank(message)
		case "iasip":
			iasip(message)
		case "remind_me":
			remind(message)
		case "care":
			care(message)
		case "pollentiek":
			pollentiek(message)
		case "presidential_order":
			presidential_order(message)
		}
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

// Returns a username that could contain markdown characters
func getUserNameUnsafe(userID int, chatID int64) (string, error) {
	user, err := getUser(userID, chatID)

	if err != nil {
		return "", err
	}

	return user.User.UserName, nil
}

// Returns a markdown escaped username
func getUserName(userID int, chatID int64) (string, error) {
	username, err := getUserNameUnsafe(userID, chatID)

	if err != nil {
		return "", err
	}

	return markdown.Escape(username), nil
}
