package models

import (
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
)

// TradeTransaction is a struct used for creating trades
type TradeTransaction struct {
	GroupID  int64        `json:"group_id" db:"group_id" binding:"required"`
	Sender   int          `json:"sender" db:"sender" binding:"required"`
	Receiver int          `json:"receiver" db:"receiver" binding:"required"`
	Amount   uint32       `json:"amount" db:"amount" binding:"required"`
	Reason   nulls.String `json:"reason"`
}

// Create maps a trade to a real transaction and creates it
func (t *TradeTransaction) Create(tx *pop.Connection) (*Transaction, error) {
	transaction := Transaction{
		GroupID:  t.GroupID,
		Sender:   nulls.NewInt(t.Sender),
		Receiver: nulls.NewInt(t.Receiver),
		Amount:   t.Amount,
		Type:     Trade,
		Cause:    t.Reason,
	}
	err := tx.Create(&transaction)

	return &transaction, err
}
