package controllers

// import (
// 	"encoding/xml"
// 	"net/http"
// 	"time"
// 	"context"
// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	apperr "github.com/sdi2200246/synaxis/internal/error"
// 	"github.com/sdi2200246/synaxis/internal/services"
// )

// type ExportBookingProvider interface {
// 	GetTicketTypesByEventID(ctx context.Context, eventID uuid.UUID) ([]services.TicketType, error)
// 	GetExportBookings(ctx context.Context, eventID uuid.UUID) ([]services.ExportBookingDetail, error)	
// }

// type ExportEventProvider interface {
// 	GetAllEvents(ctx context.Context) ([]services.DetailedEvent, error)
// }

// type ExportGeoLocation struct {
// 	Latitude  float64 `xml:"Latitude,attr" json:"latitude"`
// 	Longitude float64 `xml:"Longitude,attr" json:"longitude"`
// }

// type ExportTicketType struct {
// 	XMLName      xml.Name  `xml:"TicketType" json:"-"`
// 	TicketTypeID uuid.UUID `xml:"TicketTypeID,attr" json:"ticket_type_id"`
// 	Name         string    `xml:"Name" json:"name"`
// 	Price        float64   `xml:"Price" json:"price"`
// 	Quantity     int       `xml:"Quantity" json:"quantity"`
// 	Available    int       `xml:"Available" json:"available"`
// }

// type ExportAttendee struct {
// 	UserID uuid.UUID `xml:"UserID,attr" json:"user_id"`
// }

// type ExportBookingXML struct {
// 	XMLName         xml.Name       `xml:"Booking" json:"-"`
// 	BookingID       uuid.UUID      `xml:"BookingID,attr" json:"booking_id"`
// 	Attendee        ExportAttendee `xml:"Attendee" json:"attendee"`
// 	Time            time.Time      `xml:"Time" json:"time"`
// 	TicketTypeRef   uuid.UUID      `xml:"TicketTypeRef" json:"ticket_type_ref"`
// 	NumberOfTickets int            `xml:"NumberOfTickets" json:"number_of_tickets"`
// 	TotalCost       float64        `xml:"TotalCost" json:"total_cost"`
// 	BookingStatus   string         `xml:"BookingStatus" json:"booking_status"`
// }

// type ExportOrganizer struct {
// 	UserID uuid.UUID `xml:"UserID,attr" json:"user_id"`
// }

// type ExportEvent struct {
// 	XMLName       xml.Name           `xml:"Event" json:"-"`
// 	EventID       uuid.UUID          `xml:"EventID,attr" json:"event_id"`
// 	Title         string             `xml:"Title" json:"title"`
// 	Categories    []string           `xml:"Category" json:"categories"`
// 	EventType     string             `xml:"EventType" json:"event_type"`
// 	Venue         string             `xml:"Venue" json:"venue"`
// 	Address       string             `xml:"Address" json:"address"`
// 	City          string             `xml:"City" json:"city"`
// 	Country       string             `xml:"Country" json:"country"`
// 	GeoLocation   *ExportGeoLocation `xml:"GeoLocation,omitempty" json:"geo_location,omitempty"`
// 	StartDateTime time.Time          `xml:"StartDateTime" json:"start_datetime"`
// 	EndDateTime   time.Time          `xml:"EndDateTime" json:"end_datetime"`
// 	Capacity      int                `xml:"Capacity" json:"capacity"`
// 	TicketTypes   struct {
// 		Items []ExportTicketType `xml:"TicketType" json:"items"`
// 	} `xml:"TicketTypes" json:"ticket_types"`
// 	Bookings struct {
// 		Items []ExportBookingXML `xml:"Booking" json:"items"`
// 	} `xml:"Bookings" json:"bookings"`
// 	Organizer   ExportOrganizer `xml:"Organizer" json:"organizer"`
// 	Status      string          `xml:"Status" json:"status"`
// 	Description string          `xml:"Description" json:"description"`
// }

// type ExportEvents struct {
// 	XMLName xml.Name      `xml:"Events" json:"-"`
// 	Events  []ExportEvent `xml:"Event" json:"events"`
// }

// type AdminExportHandler struct {
// 	eventService   ExportEventProvider
// 	bookingService ExportBookingProvider
// }

// func NewAdminExportHandler(es *services.EventService, bs *services.BookingService) *AdminExportHandler {
// 	return &AdminExportHandler{eventService: es, bookingService: bs}
// }

// func (h *AdminExportHandler) Export(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	events, err := h.eventService.GetAllEvents(ctx)
// 	if err != nil {
// 		apperr.Handle(c, err)
// 		return
// 	}

// 	exportEvents := make([]ExportEvent, 0, len(events))
// 	for _, ev := range events {
// 		tickets, err := h.bookingService.GetTicketTypesByEventID(ctx, ev.ID)
// 		if err != nil {
// 			apperr.Handle(c, err)
// 			return
// 		}
// 		bookings, err := h.bookingService.GetExportBookings(ctx, ev.ID)
// 		if err != nil {
// 			apperr.Handle(c, err)
// 			return
// 		}
// 		exportEvents = append(exportEvents, buildExportEvent(ev, tickets, bookings))
// 	}

// 	payload := ExportEvents{Events: exportEvents}

// 	switch c.NegotiateFormat(gin.MIMEXML, gin.MIMEJSON) {
// 	case gin.MIMEXML:
// 		c.Header("Content-Disposition", `attachment; filename="events.xml"`)
// 		c.XML(http.StatusOK, payload)
// 	default:
// 		c.Header("Content-Disposition", `attachment; filename="events.json"`)
// 		c.JSON(http.StatusOK, payload)
// 	}
// }


