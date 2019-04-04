package models

import (
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
)

// PaymentTransaction is a struct used for creating payments
type PaymentTransaction struct {
	*baseTransaction
	GroupID int64        `json:"group_id" db:"group_id" binding:"required"`
	Sender  int          `json:"sender" db:"sender" binding:"required"`
	Cause   nulls.String `json:"cause" db:"cause" binding:"required"`
}

// Create maps a payment to a real transaction and creates it
func (p *PaymentTransaction) Create(tx *pop.Connection) (*Transaction, error) {
	var wallet Wallet

	if err := wallet.Get(tx, p.GroupID, nulls.NewInt(p.Sender)); err != nil {
		return nil, err
	}

	transaction := Transaction{
		GroupID: p.GroupID,
		Sender:  nulls.NewInt(p.Sender),
		Cause:   p.Cause,
		Type:    Payment,
	}
	err := tx.Create(&transaction)
	return &transaction, err
}
