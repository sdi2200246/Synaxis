package entities

import (
	"time"

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

func (e Event) ApproveDeletion()bool{
    return e.Status != "CANCELLED"
}

func (e Event) IsBookingAvailable()bool{
    return e.Status == "PUBLISHED"
}

func (e Event) AllowsTicketModification() bool {
    return e.Status == "DRAFT" || e.Status == "PUBLISHED"
}

func (e Event) HasCapacityFor(currentSum, additional int) bool {
    return currentSum + additional <= e.Capacity
}

func (e Event) ApprovePublication()error{
    
    if time.Now().After(e.StartDatetime){
        return  apperr.ErrCannotPublishPastEvent
    }

    if e.Status != "DRAFT"{
        return  apperr.ErrInvalidEventStatus
    }

    return nil
}

func (e Event) ApproveCancellation() error {
    if  e.Status != "PUBLISHED"{
        return  apperr.ErrInvalidEventStatus
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

