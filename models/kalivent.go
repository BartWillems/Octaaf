package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
)

type Kalivent struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UserID    int       `json:"user_id" db:"user_id"`
	Type      string    `json:"type" db:"type"`
}

type Kalivents []Kalivent

type KaliStat struct {
	Count  int    `json:"count" db:"kcount"`
	UserID int    `json:"user_id" db:"user_id"`
	Type   string `json:"type" db:"type"`
}

type KaliStats []KaliStat

func (k Kalivent) String() string {
	jk, _ := json.Marshal(k)
	return string(jk)
}

func (k Kalivents) String() string {
	jk, _ := json.Marshal(k)
	return string(jk)
}

func (k *KaliStat) TableName() string {
	return "kalivents"
}

func (k *KaliStats) TableName() string {
	return "kalivents"
}

func (k *KaliStats) Top(tx *pop.Connection, kind string) error {
	err := tx.Select("user_id", "COUNT(type) kcount").
		Where("type=?", kind).
		GroupBy("user_id").
		Order("kcount DESC").
		All(k)

	return err
}
