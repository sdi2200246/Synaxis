package entities

import "github.com/google/uuid"

type Venue struct {
    ID        uuid.UUID `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    Address   string    `json:"address" db:"address"`
    City      string    `json:"city" db:"city"`
    Country   string    `json:"country" db:"country"`
    Latitude  *float64  `json:"latitude" db:"latitude"`
    Longitude *float64  `json:"longitude" db:"longitude"`
    Capacity  *int      `json:"capacity" db:"capacity"`
}
