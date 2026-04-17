package controllers

import (
    "time"
    "github.com/google/uuid"
    "github.com/sdi2200246/synaxis/internal/services"
)

type VenueResponse struct {
    ID        uuid.UUID `json:"id"`
    Name      string    `json:"name"`
    Address   string    `json:"address"`
    City      string    `json:"city"`
    Country   string    `json:"country"`
    Latitude  *float64  `json:"latitude,omitempty"`
    Longitude *float64  `json:"longitude,omitempty"`
    Capacity *int       `json:"capacity"`    
}

type CategoryResponce struct{
    ID  uuid.UUID `json:"id"`
    Name string   `json:"name"`
}

type EventResponse struct {
    ID            uuid.UUID     `json:"id"`
    Title         string        `json:"title"`
    EventType     string        `json:"event_type"`
    Status        string        `json:"status"`
    Description   string        `json:"description"`
    Capacity      int           `json:"capacity"`
    StartDatetime time.Time     `json:"start_datetime"`
    EndDatetime   time.Time     `json:"end_datetime"`
    Venue         VenueResponse `json:"venue"`
    Categories   *[]CategoryResponce `json:"categories"`
}

type AdminUserResponse struct {
    ID        uuid.UUID `json:"id"`
    Username  string    `json:"username"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Email     string    `json:"email"`
    Address   string    `json:"address"`
    City      string    `json:"city"`
    Country   string    `json:"country"`
    TaxID     string    `json:"tax_id"`
    Status    string    `json:"status"`
    Phone     string    `json:"phone"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt *time.Time `json:"updated_at"`
}

type UserBookingResponse struct {
	ID              uuid.UUID `json:"id"`
	TicketTypeID    uuid.UUID `json:"ticket_type_id"`
	TicketName      string    `json:"ticket_name"`
	NumberOfTickets int       `json:"number_of_tickets"`
	TotalCost       float64   `json:"total_cost"`
	Status          string    `json:"status"`
	BookedAt        time.Time `json:"booked_at"`
	EventID         uuid.UUID `json:"event_id"`
	EventTitle      string    `json:"event_title"`
	EventStart      time.Time `json:"event_start"`
	VenueName       string    `json:"venue_name"`
	VenueCity       string    `json:"venue_city"`
	VenueLatitude   *float64  `json:"venue_latitude,omitempty"`
	VenueLongitude  *float64  `json:"venue_longitude,omitempty"`
}


type EventBookingResponse struct {
	ID              uuid.UUID `json:"id"`
	TicketName      string    `json:"ticket_name"`
	NumberOfTickets int       `json:"number_of_tickets"`
	TotalCost       float64   `json:"total_cost"`
	BookedAt        time.Time `json:"booked_at"`
	AttendeeName    string    `json:"attendee_name"`
	AttendeeEmail   string    `json:"attendee_email"`
	AttendeePhone   *string   `json:"attendee_phone,omitempty"`
}


func ToEventResponse(ev services.DetailedEvent) EventResponse {

    categories := make([]CategoryResponce , 0)

    for _,c:= range ev.Categories{
        categories = append(categories,  CategoryResponce{ID:c.ID,Name: c.Name})
    }

    return EventResponse{
        ID:            ev.ID,
        Title:         ev.Title,
        EventType:     ev.EventType,
        Status:        ev.Status,
        Description:   ev.Description,
        Capacity:      ev.Capacity,
        StartDatetime: ev.StartDatetime,
        EndDatetime:   ev.EndDatetime,
        Venue: VenueResponse{
            ID:        ev.VenueID,
            Name:      ev.VenueName,
            Address:   ev.VenueAddress,
            City:      ev.VenueCity,
            Country:   ev.VenueCountry,
            Latitude:  ev.VenueLatitude,
            Longitude: ev.VenueLongitude,
        },
        Categories: &categories,
    }
}

func ToEventListResponse(events []services.DetailedEvent) []EventResponse {
    result := make([]EventResponse, len(events))
    for i, ev := range events {
        result[i] = ToEventResponse(ev)
    }
    return result
}
func ToUserBookingListResponse(bookings []services.UserBookingDetail) []UserBookingResponse {
	result := make([]UserBookingResponse, len(bookings))
	for i, b := range bookings {
		result[i] = UserBookingResponse{
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
	return result
}

func ToEventBookingListResponse(bookings []services.EventBookingDetail) []EventBookingResponse {
	result := make([]EventBookingResponse, len(bookings))
	for i, b := range bookings {
		result[i] = EventBookingResponse{
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
	return result
}


func buildExportEvent(ev services.DetailedEvent,tickets []services.TicketType,bookings []services.ExportBookingDetail,) ExportEvent {
	categories := make([]string, 0, len(ev.Categories))
	for _, c := range ev.Categories {
		categories = append(categories, c.Name)
	}

	var geo *ExportGeoLocation
	if ev.VenueLatitude != nil && ev.VenueLongitude != nil {
		geo = &ExportGeoLocation{
			Latitude:  *ev.VenueLatitude,
			Longitude: *ev.VenueLongitude,
		}
	}

	exportTickets := make([]ExportTicketType, len(tickets))
	for i, t := range tickets {
		exportTickets[i] = ExportTicketType{
			TicketTypeID: t.ID,
			Name:         t.Name,
			Price:        t.Price,
			Quantity:     t.Quantity,
			Available:    t.Available,
		}
	}

	exportBookings := make([]ExportBookingXML, len(bookings))
	for i, b := range bookings {
		exportBookings[i] = ExportBookingXML{
			BookingID:       b.ID,
			Attendee:        ExportAttendee{UserID: b.AttendeeID},
			Time:            b.BookedAt,
			TicketTypeRef:   b.TicketTypeID,
			NumberOfTickets: b.NumberOfTickets,
			TotalCost:       b.TotalCost,
			BookingStatus:   b.Status,
		}
	}

	result := ExportEvent{
		EventID:       ev.ID,
		Title:         ev.Title,
		Categories:    categories,
		EventType:     ev.EventType,
		Venue:         ev.VenueName,
		Address:       ev.VenueAddress,
		City:          ev.VenueCity,
		Country:       ev.VenueCountry,
		GeoLocation:   geo,
		StartDateTime: ev.StartDatetime,
		EndDateTime:   ev.EndDatetime,
		Capacity:      ev.Capacity,
		Organizer:     ExportOrganizer{UserID: ev.OrganizerID},
		Status:        ev.Status,
		Description:   ev.Description,
	}
	result.TicketTypes.Items = exportTickets
	result.Bookings.Items = exportBookings

	return result
}