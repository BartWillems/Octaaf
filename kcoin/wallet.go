package kcoin

import (
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/opentracing/opentracing-go"
	kalicoin "gitlab.com/bartwillems/kalicoin/pkg/models"
)

type walletError struct {
	Msg string `json:"msg"`
}

// Wallet returns the wallet for a user
func Wallet(message *tgbotapi.Message, span opentracing.Span) (kalicoin.Wallet, error) {
	var wallet kalicoin.Wallet
	var wError walletError
	_, err := GetClient().
		Get("wallets/").
		// TODO: Add Jaeger headers
		Path("group/").
		Path(fmt.Sprintf("%v/", message.Chat.ID)).
		Path("owner/").
		Path(fmt.Sprintf("%v", message.From.ID)).
		Receive(&wallet, &wError)

	if err != nil {
		return kalicoin.Wallet{}, err
	}

	if (wError != walletError{}) {
		return kalicoin.Wallet{}, errors.New(wError.Msg)
	}

	return wallet, nil
}
