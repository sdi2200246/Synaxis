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

type EventCategory struct {
	ID       uuid.UUID
	Name     string
	ParentID *uuid.UUID
}

type UpdateEventInput struct{
	Title       *string
	EventType   *string
	VenueID     *uuid.UUID
	Description *string
	CategoryIDs *[]uuid.UUID
	Status 		*string
}


type Event struct {
    ID            uuid.UUID
    OrganizerID   uuid.UUID
    VenueID       uuid.UUID
    Title         string
    EventType     string
    Status        string
    Description   string
    Capacity      int
    StartDatetime time.Time
    EndDatetime   time.Time
    CreatedAt     time.Time
}

type EventFilterInput struct{
	OrganizerID   *uuid.UUID
	Status		  *string
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
    categoryProvider interfaces.CategoriesRepo
    bookingsProvider interfaces.BookingRepository
	ticketsProvider interfaces.TicketTypeRepository
}

func NewEventService(r interfaces.EventRepository ,cr interfaces.CategoriesRepo   ,br  interfaces.BookingRepository , tr interfaces.TicketTypeRepository)*EventService{
	return  &EventService{eventRepo:r , categoryProvider: cr, bookingsProvider: br , ticketsProvider: tr}
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
		Title: event.Title,
		EventType: event.EventType,
		VenueID:   event.VenueID,
		Description: event.Description,
		CategoryIDs: event.CategoryIDs,
	}
	return s.eventRepo.Update(ctx , eventID , updateEvent)
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

func (s *EventService) List(ctx context.Context, input EventFilterInput) ([]Event, bool, error) {
    filter := entities.EventFilter{
        OrganizerID: input.OrganizerID,
        Status:      input.Status,
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

    events, hasMore, err := s.eventRepo.GetbyFilter(ctx, filter)
    if err != nil {
        return nil, false, err
    }

    result := make([]Event, len(events))
    for i, e := range events {
        result[i] = toEvent(e)
    }
    return result, hasMore, nil
}

func (s *EventService) GetAllEvents(ctx context.Context) ([]Event, error) {
	events, err := s.eventRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]Event, len(events))
	for i, e := range events {
		result[i] = toEvent(e)
	}
	return result, nil
}

func toEvent(e entities.Event)Event {
	return Event{
		ID:            e.ID,
		OrganizerID:   e.OrganizerID,
		VenueID:       e.VenueID,
		Title:         e.Title,
		EventType:     e.EventType,
		Status:        e.Status,
		Description:   e.Description,
		Capacity:      e.Capacity,
		StartDatetime: e.StartDatetime,
		EndDatetime:   e.EndDatetime,
		CreatedAt:     e.CreatedAt,
	}
}

func (s *EventService) Delete(ctx context.Context, eventID uuid.UUID) error {
    event, err := s.eventRepo.GetByID(ctx, eventID)
    if err != nil {
        return err
    }

    bookingsCount , err := s.bookingsProvider.CountByEventID(ctx , eventID)
    if err != nil {
        return apperr.ErrInternal
    }

    if bookingsCount > 0 {
        return  apperr.ErrConflict
    }

    if !event.ApproveDeletion(){
        return apperr.ErrConflict
    }

    return s.eventRepo.Delete(ctx, eventID)
}


func (s *EventService) Publish(ctx context.Context , eventID uuid.UUID) error{

	event , err := s.eventRepo.GetByID(ctx , eventID)
	if err != nil{
		return err
	}
	published_tickets , err := s.ticketsProvider.SumQuantityByEventID(ctx , eventID)
	if err != nil{
		return err
	}

	if published_tickets <= 0{
		return  apperr.ErrCannotPublishWithoutTickets
	} 

	if err := event.ApprovePublication(); err != nil{
		return err
	}

	status := "PUBLISHED"
	updateEvent := entities.UpdateEvent{Status: &status}
	return  s.eventRepo.Update(ctx , eventID , updateEvent)
	
}


func (s *EventService) GetByID(ctx context.Context, id uuid.UUID) (Event, error) {
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return Event{}, err
	}
	return toEvent(event), nil
}

func (s *EventService) GetEventCategories(ctx context.Context, eventID uuid.UUID) ([]EventCategory, error) {
	categories, err := s.categoryProvider.GetByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	result := make([]EventCategory, len(categories))
	for i, c := range categories {
		result[i] = EventCategory{
			ID:       c.ID,
			Name:     c.Name,
			ParentID: c.ParentID,
		}
	}
	return result, nil
}