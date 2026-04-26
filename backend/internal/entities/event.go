package entities

import (
	"time"
    "fmt"
	"github.com/google/uuid"
	apperr "github.com/sdi2200246/synaxis/internal/error"
)

type EventStatus string

type Event struct {
    ID            uuid.UUID `db:"id"`
    OrganizerID   uuid.UUID `db:"organizer_id"`
    VenueID       uuid.UUID `db:"venue_id"`
    Title         string    `db:"title"`
    EventType     string    `db:"event_type"`
    Status        string    `db:"status"`
    Description   string    `db:"description"`
    Capacity      int       `db:"capacity"`
    StartDatetime time.Time `db:"start_datetime"`
    EndDatetime   time.Time `db:"end_datetime"`
    CreatedAt     time.Time `db:"created_at"`
}

func (e Event) ApproveDeletion() error {
    if e.Status == "CANCELLED" {
        return fmt.Errorf("cancelled events cannot be deleted: %w", apperr.ErrConflict)
    }
    return nil
}

func (e Event) IsBookingAvailable() error {
    if e.Status != "PUBLISHED" {
        return fmt.Errorf("bookings are not available for %s events: %w", e.Status, apperr.ErrConflict)
    }
    return nil
}

func (e Event) AllowsTicketModification() error {
    if e.Status != "DRAFT" && e.Status != "PUBLISHED" {
        return fmt.Errorf("cannot modify tickets for %s events: %w", e.Status, apperr.ErrConflict)
    }
    return nil
}

func (e Event) HasCapacityFor(currentSum, additional int) error {
    if currentSum+additional > e.Capacity {
        return fmt.Errorf("adding %d tickets would exceed event capacity of %d (current: %d): %w",
            additional, e.Capacity, currentSum, apperr.ErrConflict)
    }
    return nil
}

func (e Event) ApprovePublication() error {
    if time.Now().After(e.StartDatetime) {
        return fmt.Errorf("event starts at %s which is in the past: %w",
            e.StartDatetime.Format(time.RFC3339), apperr.ErrConflict)
    }
    if e.Status != "DRAFT" {
        return fmt.Errorf("cannot publish a %s event: %w", e.Status, apperr.ErrConflict)
    }
    return nil
}

func (e Event) ApproveCancellation() error {
    if e.Status != "PUBLISHED" {
        return fmt.Errorf("cannot cancel a %s event: %w", e.Status, apperr.ErrConflict)
    }
    return nil
}

type UpdateEvent struct{
    Title       *string
	EventType   *string
	VenueID     *uuid.UUID
	Description *string
	CategoryIDs *[]uuid.UUID
	Status 		*string
}

type EventFilter struct {
    OrganizerID   *uuid.UUID
	Status		  *string
    CategoryIDs   []uuid.UUID
    Title         *string
    Description   *string
    City          *string
    Country       *string
    StartAfter    *time.Time
    StartBefore   *time.Time
    MinPrice      *float64
    MaxPrice      *float64
    Limit         int
    Offset        int
}

