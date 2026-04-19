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

type CreateBookingInput struct{
	TicketTypeID uuid.UUID
	UserID		 uuid.UUID
	Quantity	 int
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


type UserBookingDetail struct {
	ID              uuid.UUID
	TicketTypeID    uuid.UUID
	TicketName      string
	NumberOfTickets int
	TotalCost       float64
	Status          string
	BookedAt        time.Time
	EventID         uuid.UUID
	EventTitle      string
	EventStart      time.Time
	VenueName       string
	VenueCity       string
	VenueLatitude   *float64
	VenueLongitude  *float64
}

type EventBookingDetail struct {
	ID              uuid.UUID
	TicketName      string
	NumberOfTickets int
	TotalCost       float64
	BookedAt        time.Time
	AttendeeName    string
	AttendeeEmail   string
	AttendeePhone   *string
}

type ExportBookingDetail struct {
	ID              uuid.UUID
	TicketTypeID    uuid.UUID
	AttendeeID      uuid.UUID
	NumberOfTickets int
	TotalCost       float64
	Status          string
	BookedAt        time.Time
}

type BookingService struct {
	ticketTypeRepo interfaces.TicketTypeRepository
	bookingRepo 	interfaces.BookingRepository
}

func NewBookingService(r interfaces.TicketTypeRepository , rb interfaces.BookingRepository) *BookingService {
	return &BookingService{ticketTypeRepo: r , bookingRepo: rb}
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


func (s *BookingService) CreateBooking(ctx context.Context, input CreateBookingInput) error {

	ticket, err := s.ticketTypeRepo.GetByID(ctx, input.TicketTypeID)
	if err != nil {
		return err
	}

	if ticket.Available < input.Quantity {
		return apperr.ErrConflict
	}
	totalCost := ticket.Price * float64(input.Quantity)

	booking := entities.Booking{
		ID:              uuid.New(),
		UserID:          input.UserID,
		TicketTypeID:    input.TicketTypeID,
		NumberOfTickets: input.Quantity,
		TotalCost:       totalCost,
		Status:          "ACTIVE",
		BookedAt:        time.Now(),
	}

	return s.bookingRepo.Create(ctx, booking)
}

func (s *BookingService) GetUserBookings(ctx context.Context, userID uuid.UUID) ([]UserBookingDetail, error) {
	bookings, err := s.bookingRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]UserBookingDetail, len(bookings))
	for i, b := range bookings {
		result[i] = UserBookingDetail{
			ID:              b.ID,
			TicketTypeID:    b.TicketTypeID,
			TicketName:      b.TicketName,
			NumberOfTickets: b.NumberOfTickets,
			TotalCost:       b.TotalCost,
			Status:          b.Status,
			BookedAt:        b.BookedAt,
			EventID:         b.EventID,
			EventTitle:      b.EventTitle,
			EventStart:      b.EventStart,
			VenueName:       b.VenueName,
			VenueCity:       b.VenueCity,
			VenueLatitude:   b.VenueLatitude,
			VenueLongitude:  b.VenueLongitude,
		}
	}
	return result, nil
}

func (s *BookingService) GetEventBookings(ctx context.Context, eventID uuid.UUID) ([]EventBookingDetail, error) {
	bookings, err := s.bookingRepo.GetByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	result := make([]EventBookingDetail, len(bookings))
	for i, b := range bookings {
		result[i] = EventBookingDetail{
			ID:              b.ID,
			TicketName:      b.TicketName,
			NumberOfTickets: b.NumberOfTickets,
			TotalCost:       b.TotalCost,
			BookedAt:        b.BookedAt,
			AttendeeName:    b.AttendeeName,
			AttendeeEmail:   b.AttendeeEmail,
			AttendeePhone:   b.AttendeePhone,
		}
	}
	return result, nil
}

func (s *BookingService) CountEventBookings(ctx context.Context, eventID uuid.UUID) (int, error) {
	return s.bookingRepo.CountByEventID(ctx, eventID)
}

func (s *BookingService) GetExportBookings(ctx context.Context, eventID uuid.UUID) ([]ExportBookingDetail, error) {
	bookings, err := s.bookingRepo.GetForExport(ctx, eventID)
	if err != nil {
		return nil, err
	}

	result := make([]ExportBookingDetail, len(bookings))
	for i, b := range bookings {
		result[i] = ExportBookingDetail{
			ID:              b.ID,
			TicketTypeID:    b.TicketTypeID,
			AttendeeID:      b.AttendeeID,
			NumberOfTickets: b.NumberOfTickets,
			TotalCost:       b.TotalCost,
			Status:          b.Status,
			BookedAt:        b.BookedAt,
		}
	}
	return result, nil
}