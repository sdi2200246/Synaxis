// backend/internal/services/ticket_type.go
package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/entities"
	"github.com/sdi2200246/synaxis/internal/interfaces"
)

type CreateTicketInput struct {
	EventID  uuid.UUID
	Name     string
	Price    float64
	Quantity int
}

type UpdateTicketTypeInput struct {
	Name     *string
	Price    *float64
	Quantity *int
}

type TicketType struct {
	ID        uuid.UUID
	EventID   uuid.UUID
	Name      string
	Price     float64
	Quantity  int
	Available int
	CreatedAt time.Time
}

type TicketTypeService struct {
	ticketTypeRepo  interfaces.TicketTypeRepository
	eventsRepo  interfaces.EventRepository
}

func NewTicketTypeService(ttr interfaces.TicketTypeRepository,er interfaces.EventRepository,) *TicketTypeService {
	return &TicketTypeService{
		ticketTypeRepo: ttr,
		eventsRepo: er,
	}
}

func (s *TicketTypeService) CreateTicketType(ctx context.Context, input CreateTicketInput) error {
	event, err := s.eventsRepo.GetByID(ctx, input.EventID)
	if err != nil {
		return err
	}
	if err = event.AllowsTicketModification() ; err != nil {
		return err
	}
	currentSum, err := s.ticketTypeRepo.SumQuantityByEventID(ctx, input.EventID)
	if err != nil {
		return err
	}

	if err = event.HasCapacityFor(currentSum, input.Quantity) ; err != nil {
    	return err
	}

	tt := entities.TicketType{
		ID:        uuid.New(),
		EventID:   input.EventID,
		Name:      input.Name,
		Price:     input.Price,
		Quantity:  input.Quantity,
		Available: input.Quantity,
		CreatedAt: time.Now(),
	}
	return s.ticketTypeRepo.Create(ctx, tt)
}

func (s *TicketTypeService) UpdateTicketType(ctx context.Context, ticketID, eventID uuid.UUID, input UpdateTicketTypeInput) error {
	event, err := s.eventsRepo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}
	if input.Quantity != nil {
		releasedTickets, err := s.ticketTypeRepo.SumQuantityByEventID(ctx, eventID)
		if err != nil {
			return err
		}
		ticketType, err := s.ticketTypeRepo.GetByID(ctx, ticketID)
		if err != nil {
			return err
		}
		if err = ticketType.CanSetQuantity(*input.Quantity); err != nil{
			return err
		}
		adjustedSum := releasedTickets - ticketType.Quantity
		if err = event.HasCapacityFor(adjustedSum, *input.Quantity) ; err != nil {
			return err 
		}	
	}

	update := entities.UpdateTicketType{
		Name:     input.Name,
		Price:    input.Price,
		Quantity: input.Quantity,
	}
	return s.ticketTypeRepo.Update(ctx, ticketID, update)
}

func (s *TicketTypeService) GetTicketTypesByEventID(ctx context.Context, eventID uuid.UUID) ([]TicketType, error) {
	tts, err := s.ticketTypeRepo.GetByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	result := make([]TicketType, len(tts))
	for i, tt := range tts {
		result[i] = toTicketType(tt)
	}
	return result, nil
}

func toTicketType(t entities.TicketType) TicketType {
	return TicketType{
		ID:        t.ID,
		EventID:   t.EventID,
		Name:      t.Name,
		Price:     t.Price,
		Quantity:  t.Quantity,
		Available: t.Available,
		CreatedAt: t.CreatedAt,
	}
}

func (s *TicketTypeService) GetByID(ctx context.Context, id uuid.UUID) (TicketType, error) {
	tt, err := s.ticketTypeRepo.GetByID(ctx, id)
	if err != nil {
		return TicketType{}, err
	}
	return toTicketType(tt), nil
}