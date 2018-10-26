package main

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	opentracing "github.com/opentracing/opentracing-go"
)

type IsMarkdown bool

// OctaafMessage extends on the tgbotapi Message
// This is so we can trace the message throughout the application
// and extend on it with functions
type OctaafMessage struct {
	*tgbotapi.Message
	Span       opentracing.Span
	IsMarkdown IsMarkdown
}

func (message *OctaafMessage) Reply(r interface{}) error {
	return message.ReplyTo(r, message.MessageID)
}

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
		_, err = Octaaf.Send(msg)
		span.SetTag("type", "text")
	case []byte:
		bytes := tgbotapi.FileBytes{Name: "image.jpg", Bytes: resp}
		msg := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)
		msg.Caption = message.CommandArguments()
		msg.ReplyToMessageID = message.MessageID
		_, err = Octaaf.Send(msg)
		span.SetTag("type", "image")
	}

	if err != nil {
		span.SetTag("error", err)
	}

	return err
}
