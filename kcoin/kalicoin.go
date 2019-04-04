package kcoin

import (
	"errors"
	"octaaf/models"

	"github.com/dghubble/sling"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gobuffalo/nulls"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	kalicoin "gitlab.com/bartwillems/kalicoin/pkg/models"
)

type Config struct {
	Enabled  bool   `toml:"enabled" env:"KALICOIN_ENABLED"`
	URI      string `toml:"uri" env:"KALICOIN_URI"`
	Username string `toml:"username" env:"KALICOIN_USERNAME"`
	Password string `toml:"password" env:"KALICOIN_PASSWORD"`
}

var config Config

// Initialize the kalicoin config
func InitConfig(c Config) {
	config = c
}

func getClient() *sling.Sling {
	return sling.New().
		Base(config.URI).
		SetBasicAuth(config.Username, config.Password)
}

// HandleTransaction returns true if a user can execute a command and returns an error if he can't
func HandleTransaction(message *tgbotapi.Message, span opentracing.Span) (bool, error) {
	if !ShouldPay(message) {
		return true, nil
	}

	transaction, err := MakePayment(message, span)

	if err != nil {
		return false, err
	}

	if transaction.Status != kalicoin.Succeeded {
		return false, errors.New(transaction.FailureReason.String)
	}

	return true, nil
}

func createTransaction(transactionKind interface{}, path string, span opentracing.Span) (kalicoin.Transaction, error) {
	var transaction kalicoin.Transaction
	_, err := getClient().
		Post("payments").
		// TODO: Add Jaeger headers
		BodyJSON(transactionKind).
		Receive(&transaction, &transaction)

	return transaction, err
}

// MakePayment sends the payment request to the kalicoin wallets
func MakePayment(message *tgbotapi.Message, span opentracing.Span) (kalicoin.Transaction, error) {
	payment := kalicoin.PaymentTransaction{
		Cause:   nulls.NewString(getPaymentType(message)),
		GroupID: message.Chat.ID,
		Sender:  message.From.ID,
	}

	transaction, err := createTransaction(payment, "payments", span)

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
