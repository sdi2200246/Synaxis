package entities

import (
    "time"
    "github.com/google/uuid"
)

type Visit struct {
    ID        uuid.UUID `json:"id" db:"id"`
    UserID    uuid.UUID `json:"user_id" db:"user_id"`
    EventID   uuid.UUID `json:"event_id" db:"event_id"`
    VisitedAt time.Time `json:"visited_at" db:"visited_at"`
}
