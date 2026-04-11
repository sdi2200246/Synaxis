package entities

import (
	"time"

	"github.com/google/uuid"
)

type TicketType struct {
    ID        uuid.UUID `db:"id"`
    EventID   uuid.UUID `db:"event_id"`
    Name      string    `db:"name"`
    Price     float64   `db:"price"`
    Quantity  int       `db:"quantity"`
    Available int       `db:"available"`
    CreatedAt time.Time `db:"created_at"`
}

type UpdateTicketType struct {
	Name     *string
	Price    *float64
	Quantity *int
}