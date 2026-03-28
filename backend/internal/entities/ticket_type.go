package entities

import (
    "github.com/google/uuid"
)

type TicketType struct {
    ID        uuid.UUID `json:"id"         db:"id"`
    EventID   uuid.UUID `json:"event_id"   db:"event_id"`
    Name      string    `json:"name"       db:"name"`
    Price     float64   `json:"price"      db:"price"`
    Quantity  int       `json:"quantity"   db:"quantity"`
    Available int       `json:"available"  db:"available"`
}