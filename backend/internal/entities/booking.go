// entities/booking.go
package entities

import (
    "time"
    "github.com/google/uuid"
)

type Booking struct {
    ID              uuid.UUID `json:"id"                db:"id"`
    UserID          uuid.UUID `json:"user_id"           db:"user_id"`
    TicketTypeID    uuid.UUID `json:"ticket_type_id"    db:"ticket_type_id"`
    NumberOfTickets int       `json:"number_of_tickets" db:"number_of_tickets"`
    TotalCost       float64   `json:"total_cost"        db:"total_cost"`
    Status          string    `json:"status"            db:"status"`
    BookedAt        time.Time `json:"booked_at"         db:"booked_at"`
}