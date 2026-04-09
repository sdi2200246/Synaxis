package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/entities"
	apperr "github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/interfaces"
)

type CreateEventInput struct {
    Title       string    
    EventType   string    
    VenueID     uuid.UUID 
    Description string    
    Capacity    int       
    StartDatetime time.Time 
    EndDatetime   time.Time 
    CategoryIDs []uuid.UUID
}

type Category struct {
    ID   uuid.UUID
    Name string
}

type UpdateEventInput struct{
	EventType   *string
	VenueID     *uuid.UUID
	Description *string
	CategoryIDs *[]uuid.UUID
	Status 		*string
}


type DetailedEvent struct {
    ID            uuid.UUID
    OrganizerID   uuid.UUID
    Title         string
    EventType     string
    Status        string
    Description   string
    Capacity      int
    StartDatetime time.Time
    EndDatetime   time.Time
    CreatedAt     time.Time

    VenueID        uuid.UUID
    VenueName      string
    VenueAddress   string
    VenueCity      string
    VenueCountry   string
    VenueLatitude  *float64
    VenueLongitude *float64
    VenueCapacity  *int

	Categories 	  []Category

}

type EventService struct{
	eventRepo interfaces.EventRepository
}

func NewEventService(r interfaces.EventRepository)*EventService{
	return  &EventService{eventRepo:r}
}

func (s*EventService)CreateEvent(ctx context.Context ,organizerID uuid.UUID , event CreateEventInput)error{

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
    err := s.eventRepo.CreateWithCategories(ctx , newEvent , event.CategoryIDs)

    if err != nil{
        return apperr.ErrInternal
    }
    return nil
}

func (s*EventService)UpdateEvent(ctx context.Context ,eventID uuid.UUID , event UpdateEventInput)error{

	updateEvent := entities.UpdateEvent{
		EventType: event.EventType,
		VenueID:   event.VenueID,
		Description: event.Description,
		CategoryIDs: event.CategoryIDs,
	}
	return s.eventRepo.Update(ctx , eventID , updateEvent)
}


func (s *EventService) GetOrganizerEvents(ctx context.Context, organizerID uuid.UUID) ([]DetailedEvent ,error) {

	events , err := s.eventRepo.GetByOrganizerID(ctx, organizerID)

	if err != nil{
		return nil , err
	}
	eventsRes := make([]DetailedEvent , 0)
	for _ , e := range(events){
		categories := make([]Category , 0)

		for _,c:= range e.Categories{
			categories = append(categories,  Category{ID:c.ID,Name: c.Name})
		}

		event:= DetailedEvent{
			ID: e.ID,
			OrganizerID: e.OrganizerID,
			Title: e.Title,
			EventType: e.EventType,
			Status: e.Status,
			Description: e.Description,
			Capacity: e.Capacity,
			StartDatetime: e.StartDatetime,
			EndDatetime: e.EndDatetime,
			CreatedAt: e.CreatedAt,

			VenueID: e.Venue.ID,
			VenueName: e.Venue.Name,
			VenueAddress: e.Venue.Address,
			VenueCity: e.Venue.City,
			VenueLatitude: e.Venue.Latitude,
			VenueLongitude: e.Venue.Longitude,
			VenueCapacity: e.Venue.Capacity,

			Categories: categories,
		}
		eventsRes = append(eventsRes, event)
	}
    return eventsRes , nil
}


