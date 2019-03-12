package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/uuid"
)

type LocationHistory struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	MessageID int       `json:"message_id" db:"message_id"`
	ChatID    int64     `json:"chat_id" db:"chat_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Lat       float64   `json:"lat" db:"lat"`
	Lng       float64   `json:"lng" db:"lng"`
	Name      string    `json:"name" db:"name"`
}

// String is not required by pop and may be deleted
func (l LocationHistory) String() string {
	jl, _ := json.Marshal(l)
	return string(jl)
}

// LocationHistories is not required by pop and may be deleted
type LocationHistories []LocationHistory

// String is not required by pop and may be deleted
func (l LocationHistories) String() string {
	jl, _ := json.Marshal(l)
	return string(jl)
}
