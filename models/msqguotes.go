package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
)

type MsgQuote struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	MessageID int       `json:"message_id" db:"message_id"`
	ChatID    int64     `json:"chat_id" db:"chat_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Type      string    `json:"type" db:"type"`
}

const (
	ImgQuote   = "img_quote"
	VodQuote   = "vod_quot"
	AudioQuote = "audio_quote"
)

// String is not required by pop and may be deleted
func (m MsgQuote) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// MsgQuotes is not required by pop and may be deleted
type MsgQuotes []MsgQuote

// String is not required by pop and may be deleted
func (m MsgQuotes) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

func (m *MsgQuote) TableName() string {
	return "msgquotes"
}

func (m *MsgQuotes) TableName() string {
	return "msgquotes"
}

// Search returns a random message quote for a given group
func (m *MsgQuote) Search(tx *pop.Connection, chatID int64, quoteType string) error {
	return tx.
		Where("chat_id = ?", chatID).
		Where("type = ?", quoteType).
		Order("random()").
		Limit(1).
		First(m)
}
