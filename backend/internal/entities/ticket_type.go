package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	apperr "github.com/sdi2200246/synaxis/internal/error"
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

func (t TicketType) HasAvailability(requested int) error {
    if t.Available < requested {
        return fmt.Errorf("only %d tickets available, requested %d : %w", t.Available, requested , apperr.ErrConflict)
    }
    return nil
}

func (t TicketType) CanSetQuantity(newQuantity int) error {
    sold := t.Quantity - t.Available
    if newQuantity < sold {
        return fmt.Errorf("cannot set quantity to %d, already sold %d tickets: %w", newQuantity, sold , apperr.ErrConflict)
    }
    return nil
}

type UpdateTicketType struct {
	Name     *string
	Price    *float64
	Quantity *int
}