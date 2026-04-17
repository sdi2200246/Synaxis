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

type EventFilterInput struct{
	CategoryIDs   []uuid.UUID
    Title         *string
    Description   *string
    City          *string
    Country       *string
    StartAfter    *time.Time
    StartBefore   *time.Time
    MinPrice      *float64
    MaxPrice      *float64
    Limit         int
    Offset        int
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

func (s *EventService) GetEventCapacity(ctx context.Context, id uuid.UUID) (int, error) {
	event , err :=  s.eventRepo.GetByID(ctx, id)
	if err != nil{
		return -1 , err
	}
	return event.Capacity , nil
}

func (s* EventService) GetEventStatus(ctx context.Context , id uuid.UUID)(string , error){
    event , err :=  s.eventRepo.GetByID(ctx, id)
	if err != nil{
		return "" , err
	}
	return event.Status , nil
}

func (s*EventService) GetEventOrganizer(ctx context.Context , id uuid.UUID)(uuid.UUID , error){
    event , err :=  s.eventRepo.GetByID(ctx, id)
	if err != nil{
		return uuid.Nil , err
	}
	return event.OrganizerID , nil
}

func (s *EventService) SearchEvents(ctx context.Context, input EventFilterInput) ([]DetailedEvent, bool, error) {
    filter := entities.EventFilter{
        CategoryIDs: input.CategoryIDs,
        Title:       input.Title,
        Description: input.Description,
        City:        input.City,
        Country:     input.Country,
        StartAfter:  input.StartAfter,
        StartBefore: input.StartBefore,
        MinPrice:    input.MinPrice,
        MaxPrice:    input.MaxPrice,
        Limit:       input.Limit,
        Offset:      input.Offset,
    }

    events, hasMore, err := s.eventRepo.SearchPublished(ctx, filter)
    if err != nil {
        return nil, false, err
    }
    result := make([]DetailedEvent, 0, len(events))
    for _, e := range events {
        result = append(result, toDetailedEvent(e))
    }

    return result, hasMore, nil
}



func toDetailedEvent(e entities.OrganizerEvent) DetailedEvent {
    categories := make([]Category, 0, len(e.Categories))
    for _, c := range e.Categories {
        categories = append(categories, Category{ID: c.ID, Name: c.Name})
    }

    return DetailedEvent{
        ID:            e.ID,
        OrganizerID:   e.OrganizerID,
        Title:         e.Title,
        EventType:     e.EventType,
        Status:        e.Status,
        Description:   e.Description,
        Capacity:      e.Capacity,
        StartDatetime: e.StartDatetime,
        EndDatetime:   e.EndDatetime,
        CreatedAt:     e.CreatedAt,

        VenueID:        e.Venue.ID,
        VenueName:      e.Venue.Name,
        VenueAddress:   e.Venue.Address,
        VenueCity:      e.Venue.City,
        VenueCountry:   e.Venue.Country,
        VenueLatitude:  e.Venue.Latitude,
        VenueLongitude: e.Venue.Longitude,
        VenueCapacity:  e.Venue.Capacity,

        Categories: categories,
    }
}

func (s *EventService) GetAllEvents(ctx context.Context) ([]DetailedEvent, error) {
	events, err := s.eventRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]DetailedEvent, 0, len(events))
	for _, e := range events {
		result = append(result, toDetailedEvent(e))
	}
	return result, nil
}