package interfaces

import (
	"context"
	"time"

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
	GetByTicketTypeID(ctx context.Context, ticketTypeID uuid.UUID) (entities.Event, error)
    Update(ctx context.Context, eventID uuid.UUID, update entities.UpdateEvent) error
	GetbyFilter(ctx context.Context, filter entities.EventFilter) ([]entities.Event, bool, error)
	GetAll(ctx context.Context) ([]entities.Event, error)
	Delete(ctx context.Context, eventID uuid.UUID) error
}


type BookingRepository interface{
	GetByTicketTypeID(ctx context.Context, ticketTypeID uuid.UUID) ([]entities.Booking, error)
	GetByID(ctx context.Context, id uuid.UUID) (entities.Booking, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]entities.Booking, error)
	GetByEventID(ctx context.Context, eventID uuid.UUID) ([]entities.Booking, error)
	GetForExport(ctx context.Context, eventID uuid.UUID) ([]entities.ExportBooking, error) 
	CountByEventID(ctx context.Context, eventID uuid.UUID) (int, error)
	Create(ctx context.Context, booking entities.Booking) error 
}

type MessagesRepository interface {
	CreateConversation(ctx context.Context,conv entities.Conversation,organizer uuid.UUID,attendee uuid.UUID,) error
	GetConversationByBookingID(ctx context.Context, bookingID uuid.UUID) (entities.Conversation, error)
	Create(ctx context.Context, msg entities.Message) error
	UpdateMessage(ctx context.Context, id uuid.UUID, mu entities.MessageUpdate) error
	GetByConversationID(ctx context.Context, conversationID uuid.UUID) ([]entities.Message, error)
	GetConversationByID(ctx context.Context, id uuid.UUID) (entities.Conversation, error) 
	GetParticipantsByConversationID(ctx context.Context,conversationID uuid.UUID,) ([]entities.ConvParticipant, error)
	GetUserConversations(ctx context.Context, userID uuid.UUID) ([]entities.Conversation, error) 
	GetParticipantsByConversationIDs(ctx context.Context,conversationIDs []uuid.UUID,) (map[uuid.UUID][]entities.ConvParticipant, error)
	GetUnreadMessagesCountByUser(ctx context.Context,userID uuid.UUID,) (map[uuid.UUID]int, error)
	GetMessagesByConversationID(ctx context.Context,conversationID uuid.UUID,) ([]entities.Message, error)
	MarkAsReadUpToMessage(ctx context.Context,conversationID uuid.UUID,userID uuid.UUID,lastMessageTime time.Time,) error 
	GetMessageByID(ctx context.Context, id uuid.UUID) (entities.Message, error)
	CreateConversationWithMessage(ctx context.Context,conv entities.Conversation,organizer uuid.UUID,attendee uuid.UUID,msg entities.Message,) error
	GetByBookingID(ctx context.Context, bookingID uuid.UUID) (entities.Conversation, error)
}

type EventBus interface{
	Publish(topic string, event any)
	Subscribe(topic string) chan any
}