package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID               uuid.UUID
	Topic            string
	Payload          json.RawMessage
	FirstProcessedAt time.Time
	LastProcessedAt  time.Time
	CompletedAt      *time.Time
}
