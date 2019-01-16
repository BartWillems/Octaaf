package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
)

type Quote struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Quote     string    `json:"quote" db:"quote"`
	UserID    int       `json:"user_id" db:"user_id"`
	ChatID    int64     `json:"chat_id" db:"chat_id"`
}

// String is not required by pop and may be deleted
func (q Quote) String() string {
	jq, _ := json.Marshal(q)
	return string(jq)
}

// Quotes is not required by pop and may be deleted
type Quotes []Quote

// String is not required by pop and may be deleted
func (q Quotes) String() string {
	jq, _ := json.Marshal(q)
	return string(jq)
}

func (q *Quote) Search(tx *pop.Connection, chatID int64, filter ...string) error {
	if len(filter) > 0 && filter[0] != "" {
		return tx.
			Where("chat_id = ?", chatID).
			Where("quote ilike '%' || ? || '%'", filter[0]).
			Order("random()").Limit(1).First(q)
	}

	return tx.Where("chat_id = ?", chatID).Order("random()").Limit(1).First(q)
}
