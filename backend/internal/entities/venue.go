package entities

import "github.com/google/uuid"

type Venue struct {
    ID        uuid.UUID `db:"id"`
    Name      string    `db:"name"`
    Address   string    `db:"address"`
    City      string    `db:"city"`
    Country   string    `db:"country"`
    Latitude  *float64  `db:"latitude"`
    Longitude *float64  `db:"longitude"`
    Capacity  *int      `db:"capacity"`
}

type VenuesFilter struct{
    Name *string
    Capacity *int
}