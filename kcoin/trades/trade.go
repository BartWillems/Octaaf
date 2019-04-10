package trades

import (
	"errors"
	"octaaf/kcoin"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gobuffalo/nulls"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	kalicoin "gitlab.com/bartwillems/kalicoin/pkg/models"
)

// MakeTrade Allows 1 user to send money to another user
func MakeTrade(message *tgbotapi.Message, span opentracing.Span) (kalicoin.Transaction, error) {

	if message.ReplyToMessage == nil || message.ReplyToMessage.From.ID == message.From.ID {
		return kalicoin.Transaction{}, errors.New("payments should be a response to other chat members")
	}

	amount, reason, err := parseAmount(message.CommandArguments())

	if err != nil {
		return kalicoin.Transaction{}, err
	}

	trade := kalicoin.TradeTransaction{
		GroupID:  message.Chat.ID,
		Sender:   message.From.ID,
		Receiver: message.ReplyToMessage.From.ID,
		Amount:   amount,
		Reason:   reason,
	}

	transaction, err := kcoin.CreateTransaction(trade, "trades", span)

	if err != nil {
		log.Errorf("Error: %v", err)
		return transaction, err
	}

	return transaction, err
}

func parseAmount(input string) (uint32, nulls.String, error) {
	if len(input) == 0 {
		return 0, nulls.NewString(""), errors.New("Invalid payment input")
	}

	s := strings.Split(input, " ")
	amount, err := strconv.ParseUint(s[0], 10, 32)

	if err != nil {
		return 0, nulls.NewString(""), errors.New("Invalid payment input")
	}

	var reason nulls.String
	if len(s) > 1 {
		reason = nulls.NewString(strings.Join(s[1:], " "))
	}

	return uint32(amount), reason, nil
}
