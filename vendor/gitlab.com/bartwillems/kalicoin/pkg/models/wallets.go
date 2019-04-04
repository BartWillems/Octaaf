package models

import (
	"errors"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	log "github.com/sirupsen/logrus"
)

// Wallet is a chat member's wallet 🤔
type Wallet struct {
	ID        int64     `json:"id" db:"id"`
	OwnerID   int       `json:"owner_id" db:"owner_id"`
	GroupID   int64     `json:"group_id" db:"group_id"`
	Capital   uint32    `json:"capital" db:"capital"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// StarterCapital is the capital that a new wallet starts with
const StarterCapital uint32 = 100

// Wallets is a collection of wallets
type Wallets []Wallet

// create will create a user's wallet with the default StarterCapital
func (w *Wallet) create(tx *pop.Connection, GroupID int64, UserID int) error {
	w.OwnerID = UserID
	w.GroupID = GroupID
	w.Capital = StarterCapital

	return tx.Save(w)
}

// Get fetches the user's wallet and creates it if it doesn't exist
func (w *Wallet) Get(tx *pop.Connection, GroupID int64, UserID nulls.Int) error {
	if !UserID.Valid {
		return errors.New("Empty UserID provided")
	}
	err := tx.Where("group_id = ?", GroupID).
		Where("owner_id = ?", UserID).
		First(w)

	if err == nil {
		return nil
	}

	log.Infof("Attempting to create the wallet for user %v in group %v with capital %v", UserID.Int, GroupID, StarterCapital)
	return w.create(tx, GroupID, UserID.Int)
}

func (w *Wallet) take(amount uint32) error {
	if w.Capital < amount {
		return errors.New("Not enough money in your wallet")
	}

	w.Capital = w.Capital - amount
	return nil
}
