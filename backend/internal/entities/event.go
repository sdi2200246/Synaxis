package entities

import (
    "time"
    "github.com/google/uuid"
)

type Event struct {
    ID            uuid.UUID `json:"id"             db:"id"`
    OrganizerID   uuid.UUID `json:"organizer_id"   db:"organizer_id"`
    VenueID       uuid.UUID `json:"venue_id"       db:"venue_id"`
    Title         string    `json:"title"          db:"title"`
    EventType     string    `json:"event_type"     db:"event_type"`
    Status        string    `json:"status"         db:"status"`
    Description   string    `json:"description"    db:"description"`
    Capacity      int       `json:"capacity"       db:"capacity"`
    StartDatetime time.Time `json:"start_datetime" db:"start_datetime"`
    EndDatetime   time.Time `json:"end_datetime"   db:"end_datetime"`
    CreatedAt     time.Time `json:"created_at"     db:"created_at"`
}

type EventWithVenue struct {
    Event
    Venue Venue `json:"venue"`
}