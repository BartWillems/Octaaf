package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	log "github.com/sirupsen/logrus"
)

type MessageCount struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Count     int       `json:"count" db:"count"`
	Diff      int       `json:"diff" db:"diff"`
}

// String is not required by pop and may be deleted
func (m MessageCount) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

func (m *MessageCount) BeforeSave(tx *pop.Connection) error {
	prevMC := MessageCount{}
	err := tx.Last(&prevMC)

	m.Diff = 0

	if err == nil && prevMC.Count > 0 {
		m.Diff = (m.Count - prevMC.Count)
		log.Debug("Previous message count: ", prevMC.Count)
	} else {
		log.Error("Unable to load previous error count: ", err)
		// This is the first message
		m.Diff = m.Count
	}
	log.Debug("Setting MessageCount with diff ", m.Diff)
	return nil
}
