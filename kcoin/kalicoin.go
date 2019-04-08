package kcoin

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/opentracing/opentracing-go"
	kalicoin "gitlab.com/bartwillems/kalicoin/pkg/models"
)

// Config is the kcoin API config
type Config struct {
	Enabled  bool   `toml:"enabled" env:"KALICOIN_ENABLED"`
	URI      string `toml:"uri" env:"KALICOIN_URI"`
	Username string `toml:"username" env:"KALICOIN_USERNAME"`
	Password string `toml:"password" env:"KALICOIN_PASSWORD"`
}

var config Config

// InitConfig initializes the kalicoin config
func InitConfig(c Config) {
	config = c
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

// CreateTransaction sends a transaction to the kalicoin
func CreateTransaction(transactionKind interface{}, path string, span opentracing.Span) (kalicoin.Transaction, error) {
	var transaction kalicoin.Transaction
	_, err := GetClient().
		Post(path).
		// TODO: Add Jaeger headers
		BodyJSON(transactionKind).
		Receive(&transaction, &transaction)

	return transaction, err
}
