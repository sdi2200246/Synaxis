package services

import(
	"github.com/sdi2200246/synaxis/internal/interfaces"
	"github.com/sdi2200246/synaxis/internal/entities"
	"context"
    "time"
    "github.com/google/uuid"
)

type CandidateEvent struct {
    Title       string    `json:"title"        binding:"required"`
    EventType   string    `json:"event_type"   binding:"required"`
    VenueID     uuid.UUID `json:"venue_id"     binding:"required"`
    Description string    `json:"description"  binding:"required"`
    Capacity    int       `json:"capacity"     binding:"required,min=1"`
    StartDatetime time.Time `json:"start_datetime" binding:"required"`
    EndDatetime   time.Time `json:"end_datetime"   binding:"required"`
    CategoryIDs []uuid.UUID `json:"category_ids" binding:"required,min=1"`
}

type Event struct {
    Title       string    `json:"title"        binding:"required"`
    EventType   string    `json:"event_type"   binding:"required"`
    Venue    	string `json:"venue_name"     binding:"required"`
    Description string    `json:"description"  binding:"required"`
	Status		string    `json:"status"  binding:"required"`
    Capacity    int       `json:"capacity"     binding:"required,min=1"`
    StartDatetime time.Time `json:"start_datetime" binding:"required"`
    EndDatetime   time.Time `json:"end_datetime"   binding:"required"`
    CategoryIDs []uuid.UUID `json:"category_ids" binding:"required,min=1"`
}

//TDO add venue repo to check capacity.
type EventService struct{
	eventRepo interfaces.EventRepository
}

func NewEventService(r interfaces.EventRepository)*EventService{
	return  &EventService{eventRepo:r}
}


func (s*EventService)CreateEvent(ctx context.Context ,organizerID uuid.UUID , event CandidateEvent)error{

	newEvent := entities.Event{
        ID:           uuid.New(),
		OrganizerID:  organizerID,
		VenueID:	  event.VenueID ,	
		Title: 		  event.Title,
		EventType: 	  event.EventType,
		Status:       "DRAFT",
		Description:  event.Description,
		Capacity: 	  event.Capacity,
		StartDatetime: event.StartDatetime,
		EndDatetime:  event.EndDatetime,
        CreatedAt:    time.Now(),
    }	

	return s.eventRepo.Create(ctx , newEvent)
}

func (s*EventService)PublishEvent(ctx context.Context ,id uuid.UUID)error{
	return s.eventRepo.Publish(ctx , id)
}

func (s*EventService)CancelEvent(ctx context.Context ,id uuid.UUID)error{
	return s.eventRepo.Cancel(ctx , id)
}