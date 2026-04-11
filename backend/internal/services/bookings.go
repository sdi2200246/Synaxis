package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/entities"
	apperr "github.com/sdi2200246/synaxis/internal/error"
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


type BookingService struct {
	ticketTypeRepo interfaces.TicketTypeRepository
}

func NewBookingService(r interfaces.TicketTypeRepository) *BookingService {
	return &BookingService{ticketTypeRepo: r}
}

func (s *BookingService) CreateTicketType(ctx context.Context, input CreateTicketInput, eventCapacity int) error {
	currentSum, err := s.ticketTypeRepo.SumQuantityByEventID(ctx, input.EventID)
	if err != nil {
		return err
	}
	if currentSum+input.Quantity > eventCapacity {
		return apperr.ErrConflict
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

func (s *BookingService) UpdateTicketType(ctx context.Context, id uuid.UUID, eventID uuid.UUID, input UpdateTicketTypeInput, eventCapacity int) error {
	if input.Quantity != nil {
		currentSum, err := s.ticketTypeRepo.SumQuantityByEventID(ctx, eventID)
		if err != nil {
			return err
		}
		existing, err := s.ticketTypeRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if currentSum-existing.Quantity+*input.Quantity > eventCapacity {
			return apperr.ErrConflict
		}
	}

	return s.ticketTypeRepo.Update(ctx, id, entities.UpdateTicketType{
		Name:     input.Name,
		Price:    input.Price,
		Quantity: input.Quantity,
	})
}

func (s *BookingService) GetTicketTypesByEventID(ctx context.Context, eventID uuid.UUID) ([]TicketType, error) {
	tickets, err := s.ticketTypeRepo.GetByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	result := make([]TicketType, len(tickets))
	for i, tt := range tickets {
		result[i] = TicketType{
			ID:        tt.ID,
			EventID:   tt.EventID,
			Name:      tt.Name,
			Price:     tt.Price,
			Quantity:  tt.Quantity,
			Available: tt.Available,
			CreatedAt: tt.CreatedAt,
		}
	}
	return result, nil
}

// func (s *BookingService) DeleteTicketType(ctx context.Context, id uuid.UUID) error {
// 	return s.ticketTypeRepo.Delete(ctx, id)
// }