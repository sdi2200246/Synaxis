package interfaces

import (
	"context"
	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/entities"
)

type TicketTypeRepository interface {
	Create(ctx context.Context, tt entities.TicketType) error
	GetByID(ctx context.Context, id uuid.UUID) (entities.TicketType, error)
	GetByEventID(ctx context.Context, eventID uuid.UUID) ([]entities.TicketType, error)
	SumQuantityByEventID(ctx context.Context, eventID uuid.UUID) (int, error)
	Update(ctx context.Context, id uuid.UUID, update entities.UpdateTicketType) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserRepository interface {
    Create(ctx context.Context, user entities.User) error
    GetByID(ctx context.Context, id uuid.UUID) (entities.User, error)
    GetByUsername(ctx context.Context, username string) (entities.User, error)
	ListUsers(ctx context.Context , filter entities.UserFilter)([]entities.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, u entities.UserUpdate) error 
}

type VenuesRepository interface {
    // GetByID(ctx context.Context, id uuid.UUID) (entities.Venue, error)
	ListVenues(ctx context.Context , filter entities.VenuesFilter) ([]entities.Venue , error)	
}

type EventRepository interface {
    CreateWithCategories(ctx context.Context, event entities.Event ,categoryIDs []uuid.UUID) error
    GetByID(ctx context.Context, id uuid.UUID) (entities.Event, error)
    GetByOrganizerID(ctx context.Context, organizerID uuid.UUID) ([]entities.OrganizerEvent, error) 
    Update(ctx context.Context, eventID uuid.UUID, update entities.UpdateEvent) error
	SearchPublished(ctx context.Context, filter entities.EventFilter) ([]entities.OrganizerEvent, bool, error)
	// Publish(ctx context.Context, id uuid.UUID) error 
	// Cancel(ctx context.Context, id uuid.UUID) error
}
