package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
)

type ImgQuote struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	MessageID int       `json:"message_id" db:"message_id"`
	ChatID    int64     `json:"chat_id" db:"chat_id"`
	UserID    int       `json:"user_id" db:"user_id"`
}

// String is not required by pop and may be deleted
func (i ImgQuote) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// ImgQuotes is not required by pop and may be deleted
type ImgQuotes []ImgQuote

// String is not required by pop and may be deleted
func (i ImgQuotes) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

func (i *ImgQuote) TableName() string {
	return "imgquotes"
}

func (i *ImgQuotes) TableName() string {
	return "imgquotes"
}

// Search returns a random image quote for a given group
func (i *ImgQuote) Search(tx *pop.Connection, chatID int64) error {
	return tx.Where("chat_id = ?", chatID).Order("random()").Limit(1).First(i)
}
