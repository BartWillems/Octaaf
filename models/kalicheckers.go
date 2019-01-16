package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
)

type Kalichecker struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UserID    int       `json:"user_id" db:"user_id"`
}

// Kalicheckers is an array of kalicheckers
type Kalicheckers []Kalichecker

// String is not required by pop and may be deleted
func (k Kalichecker) String() string {
	jk, _ := json.Marshal(k)
	return string(jk)
}

type KalicheckerStat struct {
	Count  int `json:"count" db:"kcount"`
	UserID int `json:"user_id" db:"user_id"`
}

// KalicheckerStats is an array of KalicheckerStats
type KalicheckerStats []KalicheckerStat

func (k *KalicheckerStat) TableName() string {
	return "kalicheckers"
}

func (k *Kalicheckers) TableName() string {
	return "kalicheckers"
}

func (k *KalicheckerStats) Top(tx *pop.Connection) error {
	return tx.Select("user_id", "COUNT(user_id) kcount").
		GroupBy("user_id").
		Order("kcount DESC").
		All(k)
}
