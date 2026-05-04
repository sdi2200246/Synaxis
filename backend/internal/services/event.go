package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/entities"
	apperr "github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/interfaces"
)

type EventCancelled struct {
    EventID uuid.UUID
}

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

type EventMedia struct{
	ID uuid.UUID
	Url string
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
	Media	  	  []EventMedia	
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
	venuesProvider interfaces.VenuesRepository
	mediaProvider interfaces.MediaRepository
	eventBus 		interfaces.EventBus
}

func NewEventService(
		r interfaces.EventRepository ,
		cr interfaces.CategoriesRepo   ,
		br  interfaces.BookingRepository , 
		tr interfaces.TicketTypeRepository ,
		eb interfaces.EventBus , 
		vr interfaces.VenuesRepository ,
		mr interfaces.MediaRepository,
	)*EventService{
		return  &EventService{
				eventRepo:r,
				categoryProvider: cr,
				bookingsProvider: br,
				ticketsProvider: tr,
				venuesProvider: vr,
				mediaProvider: mr,
				eventBus: eb,
			}
}
func (s*EventService)CreateEvent(ctx context.Context ,organizerID uuid.UUID , event CreateEventInput)error{

	if time.Now().After(event.StartDatetime){
		return  fmt.Errorf("Event must start after current date : %w" , apperr.ErrBadInput)
	}
	venue, err := s.venuesProvider.GetByID(ctx, event.VenueID)
	if err != nil {
		return err
	}
	if err := venue.HasCapacityFor(event.Capacity); err != nil {
		return err
	}

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
    return  s.eventRepo.CreateWithCategories(ctx , newEvent , event.CategoryIDs)
}

func (s *EventService) UpdateEvent(ctx context.Context,callerID uuid.UUID, eventID uuid.UUID, input UpdateEventInput) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if err = validateOwnership(callerID , event.OrganizerID) ; err != nil{
		return err
	}

	if input.Status != nil && *input.Status == "PUBLISHED" {
		publishedTickets , err := s.ticketsProvider.SumQuantityByEventID(ctx , eventID)
		if err != nil{
			return err
		}

		if publishedTickets <= 0 {
			return fmt.Errorf("cannot publish event with out released tickets :%w", apperr.ErrConflict)
		}

		if err := event.ApprovePublication(); err != nil{
			return err
		}

		if err = s.eventRepo.Update(ctx, eventID, entities.UpdateEvent{Status: input.Status}); err != nil {
			return err
		}

		return nil
	}

	if input.Status != nil && *input.Status == "CANCELLED" {
		if err = event.ApproveCancellation(); err != nil {
			return err
		}
		if err = s.eventRepo.Update(ctx, eventID, entities.UpdateEvent{Status: input.Status}); err != nil {
			return err
		}
		s.eventBus.Publish("EventCancelled", EventCancelled{EventID: eventID})
		return nil
	}

	return s.eventRepo.Update(ctx, eventID, entities.UpdateEvent{
		Title:       input.Title,
		EventType:   input.EventType,
		VenueID:     input.VenueID,
		Description: input.Description,
		CategoryIDs: input.CategoryIDs,
	})
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

func (s *EventService) List(ctx context.Context, callerID *uuid.UUID, input EventFilterInput) ([]Event, bool, error) {
	if callerID == nil {
		s := "PUBLISHED"
		input.Status = &s
	}

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

	if len(events) == 0 {
		return []Event{}, false, nil
	}

	eventIDs := make([]uuid.UUID, len(events))
	for i, e := range events {
		eventIDs[i] = e.ID
	}

	mediaByEvent, err := s.mediaProvider.GetByEventIDs(ctx, eventIDs)
	if err != nil {
		return nil, false, err
	}

	result := make([]Event, len(events))
	for i, e := range events {
		evt := toEvent(e)
		evt.Media = make([]EventMedia , 0)
		for _, m := range mediaByEvent[e.ID] {
			evt.Media = append(evt.Media,
				EventMedia{
					ID:m.ID,
					Url:fmt.Sprintf("/media/events/%s/%s", e.ID, m.Filename),
				})
		}
		result[i] = evt
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

func (s *EventService) Delete(ctx context.Context,callerID uuid.UUID,eventID uuid.UUID) error {
    event, err := s.eventRepo.GetByID(ctx, eventID)
    if err != nil {
        return err
    }

	if err = validateOwnership(callerID , event.OrganizerID) ; err != nil{
		return err
	}

    bookingsCount , err := s.bookingsProvider.CountByEventID(ctx , eventID)
    if err != nil {
        return apperr.ErrInternal
    }

    if bookingsCount > 0 {
    	return fmt.Errorf("cannot delete event with %d existing bookings:%w", bookingsCount , apperr.ErrConflict)
	}

    if err = event.ApproveDeletion() ; err != nil{
        return err
    }

    return s.eventRepo.Delete(ctx, eventID)
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