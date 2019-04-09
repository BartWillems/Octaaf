package models

import (
	"errors"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	log "github.com/sirupsen/logrus"
)

// TransactionType indicates the payment type
type TransactionType string

// TransactionStatus indicates the current status of a payment
type TransactionStatus string

const (
	// Trade means giving kalicoins from one chat member to another
	Trade TransactionType = "trade"
	// Payment means that a user uses his kalicoins for chat features
	Payment TransactionType = "payment"
	// Reward is money that a user receives for participating in certain bot events
	Reward TransactionType = "reward"

	// Pending indicates that the payment has not been handled
	Pending TransactionStatus = "pending"
	// Succeeded indicates a successfull payment
	Succeeded TransactionStatus = "succeeded"
	// Failed indicates an invalid payment, this could be caused by:
	// - insufficient capital
	// - invalid sender/receiver
	// - general db errors
	Failed TransactionStatus = "failed"
)

// Transaction is a log of a payment that happened
type Transaction struct {
	ID            int64             `json:"id" db:"id"`
	Type          TransactionType   `json:"type" db:"type" binding:"required"`
	Status        TransactionStatus `json:"status" db:"status"`
	GroupID       int64             `json:"group_id" db:"group_id"`
	Sender        nulls.Int         `json:"sender" db:"sender" binding:"required"`
	Receiver      nulls.Int         `json:"receiver" db:"receiver"`
	Amount        uint32            `json:"amount" db:"amount" binding:"required"`
	Cause         nulls.String      `json:"cause" db:"cause"`
	FailureReason nulls.String      `json:"failure_reason" db:"failure_reason"`
	CreatedAt     time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at" db:"updated_at"`
}

// Transactions is an array of transactions
type Transactions []Transaction

// GetByType fetches the transactions
func (t *Transactions) GetByType(tx *pop.Connection, tType TransactionType) error {
	return tx.Where("type = ?", tType).All(&t)
}

// PriceTable is a hashmap of the payment types with their prices
var PriceTable = map[TransactionType]map[string]uint32{
	Payment: {
		"roll":  2,
		"all":   10,
		"quote": 10,
	},
	Reward: {
		"checkin":  100,
		"kalivent": 20,
	},
}

// getAmount returns what a user should pay for payments
// or what he should receive for rewards
func (t *Transaction) getAmount() (uint32, error) {
	_, ok := PriceTable[t.Type]

	if !ok {
		return 0, errors.New("Invalid transaction type")
	}

	amount, ok := PriceTable[t.Type][t.Cause.String]

	if !ok {
		return 0, errors.New("Invalid transaction cause")
	}

	return amount, nil
}

// BeforeCreate ensures the correct transaction status, wallet config and payment amount
// This is needed because golang sets the status to ""
// which is not recognized as NULL in postgres
// The transactions have foreign keys to wallets, so a wallet is created
// for the transcation's initiator if the wallet is missing
func (t *Transaction) BeforeCreate(tx *pop.Connection) error {
	t.Status = Pending

	err := tx.Transaction(func(tx *pop.Connection) error {
		switch t.Type {

		case Payment:
			var wallet Wallet
			return wallet.Get(tx, t.GroupID, t.Sender)

		case Trade:
			var senderWallet, receiverWallet Wallet
			if err := senderWallet.Get(tx, t.GroupID, t.Sender); err != nil {
				return err
			}

			if err := receiverWallet.Get(tx, t.GroupID, t.Receiver); err != nil {
				return err
			}

		case Reward:
			var wallet Wallet
			return wallet.Get(tx, t.GroupID, t.Receiver)
		}

		return nil
	})

	if err != nil {
		return err
	}

	if t.Type == Payment || t.Type == Reward && t.Amount == 0 {
		amount, err := t.getAmount()

		if err != nil {
			return err
		}

		t.Amount = amount
	}

	return nil
}

// AfterCreate will trigger the correct action on a wallet
// This is either a normal payment, or a trade
func (t *Transaction) AfterCreate(tx *pop.Connection) error {
	err := tx.Transaction(func(tx *pop.Connection) error {
		var err error
		switch t.Type {

		case Trade:
			log.Infof("Trading %v kalicoins from %v to %v", t.Amount, t.Sender, t.Receiver)
			err = t.trade(tx)

		case Payment:
			log.Infof("User %v pays %v", t.Sender, t.Amount)
			err = t.pay(tx)

		case Reward:
			log.Infof("Rewarding %v with %v", t.Receiver, t.Amount)
			err = t.receive(tx)

		default:
			err = errors.New("Invalid transaction type")
		}

		if err == nil {
			t.Status = Succeeded
			err = tx.Update(t)
		}

		return err
	})

	if err != nil {
		t.Status = Failed
		t.FailureReason = nulls.NewString(err.Error())
		return tx.Update(t)
	}

	return nil
}

// Trade takes money from the transaction sender's wallet and adds it to the receiver's wallet
func (t *Transaction) trade(tx *pop.Connection) error {
	if err := t.pay(tx); err != nil {
		return err
	}

	if err := t.receive(tx); err != nil {
		return err
	}

	return nil
}

// Pay takes money from the transaction sender's wallet
func (t *Transaction) pay(tx *pop.Connection) error {
	var wallet Wallet

	if err := wallet.Get(tx, t.GroupID, t.Sender); err != nil {
		return err
	}

	if err := wallet.take(t.Amount); err != nil {
		return err
	}

	if err := tx.Update(&wallet); err != nil {
		return err
	}

	return nil
}

// receive is a transaction where a user receives money
func (t *Transaction) receive(tx *pop.Connection) error {
	var wallet Wallet

	if err := wallet.Get(tx, t.GroupID, t.Receiver); err != nil {
		return err
	}

	wallet.Capital = wallet.Capital + t.Amount

	return tx.Update(&wallet)
}
