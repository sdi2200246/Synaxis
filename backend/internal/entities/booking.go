// entities/booking.go
package entities

import (
    "time"
    "github.com/google/uuid"
)

type Booking struct {
    ID              uuid.UUID `db:"id"`
    UserID          uuid.UUID `db:"user_id"`
    TicketTypeID    uuid.UUID `db:"ticket_type_id"`
    NumberOfTickets int       `db:"number_of_tickets"`
    TotalCost       float64   `db:"total_cost"`
    Status          string    `db:"status"`
    BookedAt        time.Time `db:"booked_at"`
}

type ExportBooking struct {
	ID              uuid.UUID
	TicketTypeID    uuid.UUID
	AttendeeID      uuid.UUID
	NumberOfTickets int
	TotalCost       float64
	Status          string
	BookedAt        time.Time
}
