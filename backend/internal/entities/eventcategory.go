package entities

import "github.com/google/uuid"

type EventCategory struct {
    EventID    uuid.UUID `json:"event_id" db:"event_id"`
    CategoryID uuid.UUID `json:"category_id" db:"category_id"`
}
