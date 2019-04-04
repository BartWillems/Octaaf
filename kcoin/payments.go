package kcoin

import (
	"octaaf/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gobuffalo/nulls"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	kalicoin "gitlab.com/bartwillems/kalicoin/pkg/models"
)

// MakePayment sends the payment request to the kalicoin wallets
func MakePayment(message *tgbotapi.Message, span opentracing.Span) (kalicoin.Transaction, error) {
	payment := kalicoin.PaymentTransaction{
		Cause:   nulls.NewString(getPaymentType(message)),
		GroupID: message.Chat.ID,
		Sender:  message.From.ID,
	}

	transaction, err := CreateTransaction(payment, "payments", span)

	if err != nil {
		log.Errorf("Error: %v", err)
		return transaction, err
	}

	return transaction, err
}

// ShouldPay returns true if a command requires payment
func ShouldPay(message *tgbotapi.Message) bool {
	paymentType := getPaymentType(message)

	_, ok := kalicoin.PriceTable[kalicoin.Payment][nulls.NewString(paymentType)]

	if !ok {
		return false
	}

	return true
}

func getPaymentType(message *tgbotapi.Message) string {
	switch message.Command() {
	case "quote", "presidential_quote", models.ImgQuote, models.VodQuote, models.AudioQuote:
		if message.ReplyToMessage == nil {
			// Only pay for creating quotes, not fetching
			return ""
		}
		return "quote"
	case "roll":
		return "roll"
	case "all":
		return "all"
	}

	return ""
}
