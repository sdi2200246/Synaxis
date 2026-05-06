package controllers

import (
	"context"
	"encoding/xml"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperr "github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/services"
)

type ExportProvider interface {
	ExportByOrganizer(ctx context.Context, organizerID uuid.UUID) ([]services.ExportEvent, error)
}

type ExportGeoLocation struct {
	Latitude  float64 `xml:"Latitude,attr" json:"latitude"`
	Longitude float64 `xml:"Longitude,attr" json:"longitude"`
}

type ExportTicketType struct {
	XMLName      xml.Name  `xml:"TicketType" json:"-"`
	TicketTypeID uuid.UUID `xml:"TicketTypeID,attr" json:"ticket_type_id"`
	Name         string    `xml:"Name" json:"name"`
	Price        float64   `xml:"Price" json:"price"`
	Quantity     int       `xml:"Quantity" json:"quantity"`
	Available    int       `xml:"Available" json:"available"`
}

type ExportAttendee struct {
	UserID uuid.UUID `xml:"UserID,attr" json:"user_id"`
}

type ExportOrganizer struct {
	UserID uuid.UUID `xml:"UserID,attr" json:"user_id"`
}

type ExportBookingXML struct {
	XMLName         xml.Name       `xml:"Booking" json:"-"`
	BookingID       uuid.UUID      `xml:"BookingID,attr" json:"booking_id"`
	Attendee        ExportAttendee `xml:"Attendee" json:"attendee"`
	Time            time.Time      `xml:"Time" json:"time"`
	TicketTypeRef   uuid.UUID      `xml:"TicketTypeRef" json:"ticket_type_ref"`
	NumberOfTickets int            `xml:"NumberOfTickets" json:"number_of_tickets"`
	TotalCost       float64        `xml:"TotalCost" json:"total_cost"`
	BookingStatus   string         `xml:"BookingStatus" json:"booking_status"`
}

type ExportEvent struct {
	XMLName       xml.Name           `xml:"Event" json:"-"`
	EventID       uuid.UUID          `xml:"EventID,attr" json:"event_id"`
	Title         string             `xml:"Title" json:"title"`
	Categories    []string           `xml:"Category" json:"categories"`
	EventType     string             `xml:"EventType" json:"event_type"`
	Venue         string             `xml:"Venue" json:"venue"`
	Address       string             `xml:"Address" json:"address"`
	City          string             `xml:"City" json:"city"`
	Country       string             `xml:"Country" json:"country"`
	GeoLocation   *ExportGeoLocation `xml:"GeoLocation,omitempty" json:"geo_location,omitempty"`
	StartDateTime time.Time          `xml:"StartDateTime" json:"start_datetime"`
	EndDateTime   time.Time          `xml:"EndDateTime" json:"end_datetime"`
	Capacity      int                `xml:"Capacity" json:"capacity"`
	TicketTypes   struct {
		Items []ExportTicketType `xml:"TicketType" json:"items"`
	} `xml:"TicketTypes" json:"ticket_types"`
	Bookings struct {
		Items []ExportBookingXML `xml:"Booking" json:"items"`
	} `xml:"Bookings" json:"bookings"`
	Organizer   ExportOrganizer `xml:"Organizer" json:"organizer"`
	Status      string          `xml:"Status" json:"status"`
	Description string          `xml:"Description" json:"description"`
}

type ExportEvents struct {
	XMLName xml.Name      `xml:"Events" json:"-"`
	Events  []ExportEvent `xml:"Event" json:"events"`
}


func ToExportEventsResponse(src []services.ExportEvent) ExportEvents {
	out := make([]ExportEvent, len(src))
	for i, ev := range src {
		out[i] = ToExportEventResponse(ev)
	}
	return ExportEvents{Events: out}
}

func ToExportEventResponse(ev services.ExportEvent) ExportEvent {
	var geo *ExportGeoLocation
	if ev.GeoLocation != nil {
		geo = &ExportGeoLocation{
			Latitude:  ev.GeoLocation.Latitude,
			Longitude: ev.GeoLocation.Longitude,
		}
	}

	tickets := make([]ExportTicketType, len(ev.TicketTypes))
	for i, t := range ev.TicketTypes {
		tickets[i] = ExportTicketType{
			TicketTypeID: t.ID,
			Name:         t.Name,
			Price:        t.Price,
			Quantity:     t.Quantity,
			Available:    t.Available,
		}
	}

	bookings := make([]ExportBookingXML, len(ev.Bookings))
	for i, b := range ev.Bookings {
		bookings[i] = ExportBookingXML{
			BookingID:       b.ID,
			Attendee:        ExportAttendee{UserID: b.AttendeeID},
			Time:            b.BookedAt,
			TicketTypeRef:   b.TicketTypeID,
			NumberOfTickets: b.NumberOfTickets,
			TotalCost:       b.TotalCost,
			BookingStatus:   b.Status,
		}
	}

	out := ExportEvent{
		EventID:       ev.ID,
		Title:         ev.Title,
		Categories:    ev.Categories,
		EventType:     ev.EventType,
		Venue:         ev.VenueName,
		Address:       ev.Address,
		City:          ev.City,
		Country:       ev.Country,
		GeoLocation:   geo,
		StartDateTime: ev.StartDatetime,
		EndDateTime:   ev.EndDatetime,
		Capacity:      ev.Capacity,
		Organizer:     ExportOrganizer{UserID: ev.OrganizerID},
		Status:        ev.Status,
		Description:   ev.Description,
	}
	out.TicketTypes.Items = tickets
	out.Bookings.Items = bookings
	return out
}

type AdminExportHandler struct {
	exportService ExportProvider
}

func NewAdminExportHandler(es *services.ExportService) *AdminExportHandler {
	return &AdminExportHandler{exportService: es}
}

func (h *AdminExportHandler) ExportByOrganizer(c *gin.Context) {
	organizerID, err := uuid.Parse(c.Query("organizer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organizer_id"})
		return
	}

	events, err := h.exportService.ExportByOrganizer(c.Request.Context(), organizerID)
	if err != nil {
		apperr.Handle(c, err)
		return
	}

	payload := ToExportEventsResponse(events)

	switch c.NegotiateFormat(gin.MIMEXML, gin.MIMEJSON) {
	case gin.MIMEXML:
		c.Header("Content-Disposition", `attachment; filename="events.xml"`)
		c.XML(http.StatusOK, payload)
	default:
		c.Header("Content-Disposition", `attachment; filename="events.json"`)
		c.JSON(http.StatusOK, payload)
	}
}