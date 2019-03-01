package main

import (
	"errors"
	"octaaf/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

const (
	development = "development"
	production  = "production"
)

// OctaafMessage extends on the tgbotapi Message
// This is so we can trace the message throughout the application
// and extend on it with functions
type OctaafMessage struct {
	*tgbotapi.Message
	Span           opentracing.Span
	IsMarkdown     bool
	KeyboardCloser bool // When true, close an open keyboard
}

// Reply to the current message
func (message *OctaafMessage) Reply(r interface{}) error {
	return message.ReplyTo(r, message.MessageID)
}

// ReplyTo sends a reply to a specific message ID
func (message *OctaafMessage) ReplyTo(r interface{}, messageID int) error {

	span := message.Span.Tracer().StartSpan(
		"reply",
		opentracing.ChildOf(message.Span.Context()),
	)
	defer span.Finish()

	var err error
	switch resp := r.(type) {
	default:
		err = errors.New("Unkown response type given")
		span.SetTag("type", "unknown")
	case string:
		msg := tgbotapi.NewMessage(message.Chat.ID, resp)
		msg.ReplyToMessageID = messageID
		if message.IsMarkdown {
			msg.ParseMode = "markdown"
		}

		if message.KeyboardCloser {
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}
		_, err = Octaaf.Send(msg)
		span.SetTag("type", "text")
	case []byte:
		bytes := tgbotapi.FileBytes{Name: "image.jpg", Bytes: resp}
		msg := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)
		msg.Caption = message.CommandArguments()
		msg.ReplyToMessageID = message.MessageID
		_, err = Octaaf.Send(msg)
		span.SetTag("type", "image")
	case models.MsgQuote:
		msg := tgbotapi.NewForward(resp.ChatID, resp.ChatID, resp.MessageID)
		msg.ReplyToMessageID = message.MessageID
		_, err = Octaaf.Send(msg)
		span.SetTag("type", "forward")
	}

	if err != nil {
		span.SetTag("error", err)
	}

	return err
}

// LogError logs an error and reports it to jaeger
func (message *OctaafMessage) LogError(err string) {
	log.Error(err)
	message.Span.SetTag("error", true)
	message.Span.SetBaggageItem("error", err)
}
