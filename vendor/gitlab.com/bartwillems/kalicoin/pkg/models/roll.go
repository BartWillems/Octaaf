package models

import (
	"errors"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
)

// RollReward is a struct used for creating roll rewards
type RollReward struct {
	GroupID    int64        `json:"group_id" binding:"required"`
	Receiver   int          `json:"receiver" binding:"required"`
	Multiplier nulls.UInt32 `json:"multiplier" binding:"required"`
}

// Create maps a reward to a real transaction and creates it
func (r *RollReward) Create(tx *pop.Connection) (*Transaction, error) {
	var transaction Transaction
	cause, err := r.GetCause()

	if err != nil {
		return &transaction, err
	}

	transaction = Transaction{
		GroupID:  r.GroupID,
		Receiver: nulls.NewInt(r.Receiver),
		Cause:    nulls.NewString("roll - " + cause),
		Amount:   r.GetAmount(),
		Type:     Reward,
	}
	err = tx.Create(&transaction)

	return &transaction, err
}

// GetAmount returns the amount of kalicoins a user should receive for his /roll attempt
func (r *RollReward) GetAmount() uint32 {
	return (r.Multiplier.UInt32 + 1) * 10 * PriceTable[Payment]["roll"]
}

// GetCause returns the name for a certain multiplier
// eg dubs for 0
func (r *RollReward) GetCause() (string, error) {
	if r.Multiplier.UInt32 > 8 {
		return "", errors.New("Invalid Multiplier, must be within ]0-8[")
	}

	names := [9]string{"dubs", "trips", "quads", "penta", "hexa", "septa", "octa", "el niño"}

	return names[r.Multiplier.UInt32], nil
}
