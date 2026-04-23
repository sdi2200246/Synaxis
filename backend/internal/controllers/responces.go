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

type CategoryResponse struct{
    ID  uuid.UUID `json:"id"`
    Name string   `json:"name"`
}

type EventResponse struct {
	ID            uuid.UUID `json:"id"`
	OrganizerID   uuid.UUID `json:"organizer_id"`
	VenueID       uuid.UUID `json:"venue_id"`
	Title         string    `json:"title"`
	EventType     string    `json:"event_type"`
	Status        string    `json:"status"`
	Description   string    `json:"description"`
	Capacity      int       `json:"capacity"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
	CreatedAt     time.Time `json:"created_at"`
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


type ConversationResponse struct {
	ID          uuid.UUID `json:"id"`
	BookingID   uuid.UUID `json:"booking_id"`
	CreatedAt   time.Time `json:"created_at"`
	UnseenCount int       `json:"unseen_count"`
}

type ConvParticipantResponse struct {
	Role   string    `json:"role"`
	UserID uuid.UUID `json:"user_id"`
}

type ConversationWithParticipantsResponse struct {
	Conversation ConversationResponse       `json:"conversation"`
	Participants []ConvParticipantResponse   `json:"participants"`
}

type MessageResponse struct {
	ID             uuid.UUID  `json:"id"`
	ConversationID uuid.UUID  `json:"conversation_id"`
	SenderID       uuid.UUID  `json:"sender_id"`
	Content        string     `json:"content"`
	IsRead         bool       `json:"is_read"`
	IsDeleted      bool       `json:"is_deleted"`
	SentAt         time.Time  `json:"sent_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

type CreateConversationResponse struct {
	ConversationID uuid.UUID `json:"conversation_id"`
}


func ToVenueResponse(v services.DetailedVenue) VenueResponse {
	return VenueResponse{
		ID:        v.ID,
		Name:      v.Name,
		Address:   v.Address,
		City:      v.City,
		Country:   v.Country,
		Latitude:  v.Latitude,
		Longitude: v.Longitude,
		Capacity:  v.Capacity,
	}
}


func ToEventResponse(ev services.Event) EventResponse {
	return EventResponse{
		ID:            ev.ID,
		OrganizerID:   ev.OrganizerID,
		VenueID:       ev.VenueID,
		Title:         ev.Title,
		EventType:     ev.EventType,
		Status:        ev.Status,
		Description:   ev.Description,
		Capacity:      ev.Capacity,
		StartDatetime: ev.StartDatetime,
		EndDatetime:   ev.EndDatetime,
		CreatedAt:     ev.CreatedAt,
	}
}

func ToEventListResponse(events []services.Event) []EventResponse {
    result := make([]EventResponse, len(events))
    for i, ev := range events {
        result[i] = ToEventResponse(ev)
    }
    return result
}
type BookingResponse struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	TicketTypeID    uuid.UUID `json:"ticket_type_id"`
	NumberOfTickets int       `json:"number_of_tickets"`
	TotalCost       float64   `json:"total_cost"`
	Status          string    `json:"status"`
	BookedAt        time.Time `json:"booked_at"`
}

func ToBookingResponse(b services.Booking) BookingResponse {
	return BookingResponse{
		ID:              b.ID,
		UserID:          b.UserID,
		TicketTypeID:    b.TicketTypeID,
		NumberOfTickets: b.NumberOfTickets,
		TotalCost:       b.TotalCost,
		Status:          b.Status,
		BookedAt:        b.BookedAt,
	}
}

func ToBookingListResponse(bookings []services.Booking) []BookingResponse {
	result := make([]BookingResponse, len(bookings))
	for i, b := range bookings {
		result[i] = ToBookingResponse(b)
	}
	return result
}

func ToCategoryListResponse(categories []services.EventCategory) []CategoryResponse {
	result := make([]CategoryResponse, len(categories))
	for i, c := range categories {
		result[i] = CategoryResponse{
			ID:       c.ID,
			Name:     c.Name,
		}
	}
	return result
}

type PublicUserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
}

func ToPublicUserResponse(u services.PublicUser) PublicUserResponse {
	return PublicUserResponse{
		ID:        u.ID,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone:     u.Phone,
	}
}


func ToConversationResponse(c services.Conversation) ConversationResponse {
	return ConversationResponse{
		ID:          c.ID,
		BookingID:   c.BookingID,
		CreatedAt:   c.CreatedAt,
		UnseenCount: c.UnseenCount,
	}
}

func ToConversationListResponse(convs []services.Conversation) []ConversationResponse {
	result := make([]ConversationResponse, len(convs))
	for i, c := range convs {
		result[i] = ToConversationResponse(c)
	}
	return result
}

func ToConvParticipantResponse(p services.ConvParticipant) ConvParticipantResponse {
	return ConvParticipantResponse{
		Role:   p.Role,
		UserID: p.UserID,
	}
}

func ToConvParticipantsResponse(ps []services.ConvParticipant) []ConvParticipantResponse {
	result := make([]ConvParticipantResponse, len(ps))
	for i, p := range ps {
		result[i] = ToConvParticipantResponse(p)
	}
	return result
}

func ToMessageResponse(m services.Message) MessageResponse {
	return MessageResponse{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		SenderID:       m.SenderID,
		Content:        m.Content,
		IsRead:         m.IsRead,
		IsDeleted:      m.Deleted,
		SentAt:         m.SentAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func ToMessageListResponse(messages []services.Message) []MessageResponse {
	result := make([]MessageResponse, len(messages))
	for i, m := range messages {
		result[i] = ToMessageResponse(m)
	}
	return result
}




// func buildExportEvent(ev services.DetailedEvent,tickets []services.TicketType,bookings []services.ExportBookingDetail,) ExportEvent {
// 	categories := make([]string, 0, len(ev.Categories))
// 	for _, c := range ev.Categories {
// 		categories = append(categories, c.Name)
// 	}

// 	var geo *ExportGeoLocation
// 	if ev.VenueLatitude != nil && ev.VenueLongitude != nil {
// 		geo = &ExportGeoLocation{
// 			Latitude:  *ev.VenueLatitude,
// 			Longitude: *ev.VenueLongitude,
// 		}
// 	}

// 	exportTickets := make([]ExportTicketType, len(tickets))
// 	for i, t := range tickets {
// 		exportTickets[i] = ExportTicketType{
// 			TicketTypeID: t.ID,
// 			Name:         t.Name,
// 			Price:        t.Price,
// 			Quantity:     t.Quantity,
// 			Available:    t.Available,
// 		}
// 	}

// 	exportBookings := make([]ExportBookingXML, len(bookings))
// 	for i, b := range bookings {
// 		exportBookings[i] = ExportBookingXML{
// 			BookingID:       b.ID,
// 			Attendee:        ExportAttendee{UserID: b.AttendeeID},
// 			Time:            b.BookedAt,
// 			TicketTypeRef:   b.TicketTypeID,
// 			NumberOfTickets: b.NumberOfTickets,
// 			TotalCost:       b.TotalCost,
// 			BookingStatus:   b.Status,
// 		}
// 	}

// 	result := ExportEvent{
// 		EventID:       ev.ID,
// 		Title:         ev.Title,
// 		Categories:    categories,
// 		EventType:     ev.EventType,
// 		Venue:         ev.VenueName,
// 		Address:       ev.VenueAddress,
// 		City:          ev.VenueCity,
// 		Country:       ev.VenueCountry,
// 		GeoLocation:   geo,
// 		StartDateTime: ev.StartDatetime,
// 		EndDateTime:   ev.EndDatetime,
// 		Capacity:      ev.Capacity,
// 		Organizer:     ExportOrganizer{UserID: ev.OrganizerID},
// 		Status:        ev.Status,
// 		Description:   ev.Description,
// 	}
// 	result.TicketTypes.Items = exportTickets
// 	result.Bookings.Items = exportBookings

// 	return result
// }