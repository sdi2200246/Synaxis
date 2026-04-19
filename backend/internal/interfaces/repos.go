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
	GetByID(ctx context.Context, id uuid.UUID) (entities.Venue, error)
	ListVenues(ctx context.Context , filter entities.VenuesFilter) ([]entities.Venue , error)	
}


type CategoriesRepo interface{
	GetByEventID(ctx context.Context, eventID uuid.UUID) ([]entities.Category, error)
} 

type EventRepository interface {
    CreateWithCategories(ctx context.Context, event entities.Event ,categoryIDs []uuid.UUID) error
    GetByID(ctx context.Context, id uuid.UUID) (entities.Event, error)
    Update(ctx context.Context, eventID uuid.UUID, update entities.UpdateEvent) error
	GetbyFilter(ctx context.Context, filter entities.EventFilter) ([]entities.Event, bool, error)
	GetAll(ctx context.Context) ([]entities.Event, error)
	Delete(ctx context.Context, eventID uuid.UUID) error
	// Publish(ctx context.Context, id uuid.UUID) error 
	// Cancel(ctx context.Context, id uuid.UUID) error
}


type BookingRepository interface{
	GetByTicketTypeID(ctx context.Context, ticketTypeID uuid.UUID) ([]entities.Booking, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]entities.UserBooking, error)
	GetByEventID(ctx context.Context, eventID uuid.UUID) ([]entities.EventBooking, error)
	GetForExport(ctx context.Context, eventID uuid.UUID) ([]entities.ExportBooking, error) 
	CountByEventID(ctx context.Context, eventID uuid.UUID) (int, error)
	Create(ctx context.Context, booking entities.Booking) error 
}