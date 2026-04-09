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