package entities

import (
	"fmt"

	"github.com/google/uuid"
	apperr "github.com/sdi2200246/synaxis/internal/error"
)

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

func (v Venue) HasCapacityFor(requested int) error {
    if v.Capacity == nil {
        return nil // no capacity limit set on venue
    }
    if requested > *v.Capacity {
        return fmt.Errorf("requested capacity %d exceeds venue capacity of %d: %w",
            requested, *v.Capacity, apperr.ErrBadInput)
    }
    return nil
}

type VenuesFilter struct{
    Name *string
    Capacity *int
}