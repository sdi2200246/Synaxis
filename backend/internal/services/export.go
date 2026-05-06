package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/entities"
	"github.com/sdi2200246/synaxis/internal/interfaces"
)

type ExportGeoLocation struct {
	Latitude  float64
	Longitude float64
}

type ExportTicketType struct {
	ID        uuid.UUID
	Name      string
	Price     float64
	Quantity  int
	Available int
}

type ExportBookingItem struct {
	ID              uuid.UUID
	AttendeeID      uuid.UUID
	TicketTypeID    uuid.UUID
	BookedAt        time.Time
	NumberOfTickets int
	TotalCost       float64
	Status          string
}

type ExportEvent struct {
	ID            uuid.UUID
	Title         string
	Categories    []string
	EventType     string
	VenueName     string
	Address       string
	City          string
	Country       string
	GeoLocation   *ExportGeoLocation
	StartDatetime time.Time
	EndDatetime   time.Time
	Capacity      int
	TicketTypes   []ExportTicketType
	Bookings      []ExportBookingItem
	OrganizerID   uuid.UUID
	Status        string
	Description   string
}

type ExportService struct {
	eventRepo    interfaces.EventRepository
	venueRepo    interfaces.VenuesRepository
	categoryRepo interfaces.CategoriesRepo
	ticketRepo   interfaces.TicketTypeRepository
	bookingRepo  interfaces.BookingRepository
}

func NewExportService(
	er interfaces.EventRepository,
	vr interfaces.VenuesRepository,
	cr interfaces.CategoriesRepo,
	tr interfaces.TicketTypeRepository,
	br interfaces.BookingRepository,
) *ExportService {
	return &ExportService{
		eventRepo:    er,
		venueRepo:    vr,
		categoryRepo: cr,
		ticketRepo:   tr,
		bookingRepo:  br,
	}
}

func (s *ExportService) ExportByOrganizer(ctx context.Context, organizerID uuid.UUID) ([]ExportEvent, error) {
	filter := entities.EventFilter{OrganizerID: &organizerID}

	events, _, err := s.eventRepo.GetbyFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	out := make([]ExportEvent, 0, len(events))
	for _, ev := range events {
		built, err := s.buildExportEvent(ctx, ev)
		if err != nil {
			return nil, err
		}
		out = append(out, built)
	}
	return out, nil
}

func (s *ExportService) buildExportEvent(ctx context.Context, ev entities.Event) (ExportEvent, error) {
	venue, err := s.venueRepo.GetByID(ctx, ev.VenueID)
	if err != nil {
		return ExportEvent{}, err
	}

	categories, err := s.categoryRepo.GetByEventID(ctx, ev.ID)
	if err != nil {
		return ExportEvent{}, err
	}

	tickets, err := s.ticketRepo.GetByEventID(ctx, ev.ID)
	if err != nil {
		return ExportEvent{}, err
	}

	bookings, err := s.bookingRepo.GetForExport(ctx, ev.ID)
	if err != nil {
		return ExportEvent{}, err
	}

	categoryNames := make([]string, len(categories))
	for i, c := range categories {
		categoryNames[i] = c.Name
	}

	var geo *ExportGeoLocation
	if venue.Latitude != nil && venue.Longitude != nil {
		geo = &ExportGeoLocation{
			Latitude:  *venue.Latitude,
			Longitude: *venue.Longitude,
		}
	}

	exportTickets := make([]ExportTicketType, len(tickets))
	for i, t := range tickets {
		exportTickets[i] = ExportTicketType{
			ID:        t.ID,
			Name:      t.Name,
			Price:     t.Price,
			Quantity:  t.Quantity,
			Available: t.Available,
		}
	}

	exportBookings := make([]ExportBookingItem, len(bookings))
	for i, b := range bookings {
		exportBookings[i] = ExportBookingItem{
			ID:              b.ID,
			AttendeeID:      b.AttendeeID,
			TicketTypeID:    b.TicketTypeID,
			BookedAt:        b.BookedAt,
			NumberOfTickets: b.NumberOfTickets,
			TotalCost:       b.TotalCost,
			Status:          b.Status,
		}
	}

	return ExportEvent{
		ID:            ev.ID,
		Title:         ev.Title,
		Categories:    categoryNames,
		EventType:     ev.EventType,
		VenueName:     venue.Name,
		Address:       venue.Address,
		City:          venue.City,
		Country:       venue.Country,
		GeoLocation:   geo,
		StartDatetime: ev.StartDatetime,
		EndDatetime:   ev.EndDatetime,
		Capacity:      ev.Capacity,
		TicketTypes:   exportTickets,
		Bookings:      exportBookings,
		OrganizerID:   ev.OrganizerID,
		Status:        ev.Status,
		Description:   ev.Description,
	}, nil
}