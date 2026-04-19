package entities

import (
    "time"
    "github.com/google/uuid"
)

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

func (e Event) ApproveDeletion()bool{
    return e.Status != "CANCELLED"
}

type OrganizerEvent struct {
    Event
    Venue Venue
    Categories []Category
}

type UpdateEvent struct{
	EventType   *string
	VenueID     *uuid.UUID
	Description *string
	CategoryIDs *[]uuid.UUID
	Status 		*string
}

type EventFilter struct {
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

