package rewards

import (
	"octaaf/kcoin"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gobuffalo/nulls"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	kalicoin "gitlab.com/bartwillems/kalicoin/pkg/models"
)

// RewardRoll attempts to pay a user when he has a good roll
func RewardRoll(message *tgbotapi.Message, span opentracing.Span, multiplier int8) (kalicoin.Transaction, error) {
	reward := kalicoin.RollReward{
		GroupID:    message.Chat.ID,
		Receiver:   message.From.ID,
		Multiplier: nulls.NewUInt32(uint32(multiplier)),
	}

	transaction, err := kcoin.CreateTransaction(reward, "rewards/roll", span)

	if err != nil {
		log.Errorf("Error: %v", err)
		return transaction, err
	}

	return transaction, nil
}
