package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/entities"
	apperr "github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/interfaces"
)

type CreateBookingInput struct{
	TicketTypeID uuid.UUID
	UserID		 uuid.UUID
	Quantity	 int
}

type Booking struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	TicketTypeID    uuid.UUID
	NumberOfTickets int
	TotalCost       float64
	Status          string
	BookedAt        time.Time
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
	ticketTypeRepo 	interfaces.TicketTypeRepository
	bookingRepo 	interfaces.BookingRepository
	eventsRepo 		interfaces.EventRepository
}

func NewBookingService(r interfaces.TicketTypeRepository , rb interfaces.BookingRepository , er interfaces.EventRepository) *BookingService {
	return &BookingService{ticketTypeRepo: r , bookingRepo: rb , eventsRepo: er}
}

func (s *BookingService) CreateBooking(ctx context.Context, input CreateBookingInput) error {

	ticket, err := s.ticketTypeRepo.GetByID(ctx, input.TicketTypeID)
	if err != nil {
		return err
	}
	event, err := s.eventsRepo.GetByID(ctx, ticket.EventID)
    if err != nil {
        return err
    }

	if !event.IsBookingAvailable(){
		return apperr.ErrConflict
	}

	if !ticket.HasAvailability(input.Quantity) {
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

func (s *BookingService) GetUserBookings(ctx context.Context, userID uuid.UUID) ([]Booking, error) {
	bookings, err := s.bookingRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]Booking, len(bookings))
	for i, b := range bookings {
		result[i] = toBooking(b)
	}
	return result, nil
}

func (s *BookingService) GetEventBookings(ctx context.Context, eventID uuid.UUID) ([]Booking, error) {
	bookings, err := s.bookingRepo.GetByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	result := make([]Booking, len(bookings))
	for i, b := range bookings {
		result[i] = toBooking(b)
	}
	return result, nil
}

func toBooking(b entities.Booking) Booking {
	return Booking{
		ID:              b.ID,
		UserID:          b.UserID,
		TicketTypeID:    b.TicketTypeID,
		NumberOfTickets: b.NumberOfTickets,
		TotalCost:       b.TotalCost,
		Status:          b.Status,
		BookedAt:        b.BookedAt,
	}
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