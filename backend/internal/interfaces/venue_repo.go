package interfaces

import (
	"context"
	"github.com/sdi2200246/synaxis/internal/entities"
)

type VenuesRepository interface {
    // GetByID(ctx context.Context, id uuid.UUID) (entities.Venue, error)
	ListVenues(ctx context.Context , filter entities.VenuesFilter) ([]entities.Venue , error)	
}