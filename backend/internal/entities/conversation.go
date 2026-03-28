package entities

import (
    "time"
    "github.com/google/uuid"
)

type Conversation struct {
    ID        uuid.UUID `json:"id" db:"id"`
    BookingID uuid.UUID `json:"booking_id" db:"booking_id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
