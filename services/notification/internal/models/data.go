package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type CreateEventData struct {
	ID      uuid.UUID
	Topic   string
	Payload json.RawMessage
}
