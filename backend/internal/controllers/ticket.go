package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	apperr "github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/interfaces"
	"github.com/sdi2200246/synaxis/internal/services"
)

type CreateTicketTypeRequest struct {
	Name     string  `json:"name"     binding:"required"`
	Price    float64 `json:"price"    binding:"required,min=0"`
	Quantity int     `json:"quantity" binding:"required,min=1"`
}

type UpdateTicketTypeRequest struct {
	Name     *string  `json:"name,omitempty"`
	Price    *float64 `json:"price,omitempty"`
	Quantity *int     `json:"quantity,omitempty"`
}

type TicketTypeResponse struct {
	ID        uuid.UUID `json:"id"`
	EventID   uuid.UUID `json:"event_id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	Available int       `json:"available"`
}

type TicketTypeHandler struct {
	bookingService   *services.BookingService
	eventsProvider interfaces.EventsProvider
}

func NewTicketTypeHandler(bs *services.BookingService, cp interfaces.EventsProvider) *TicketTypeHandler {
	return &TicketTypeHandler{bookingService: bs, eventsProvider: cp}
}

func (h *TicketTypeHandler) Create(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid event id"})
		return
	}

	var input CreateTicketTypeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "invalid input", "details": err.Error()})
		return
	}

	capacity, err := h.eventsProvider.GetEventCapacity(c.Request.Context(), eventID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	err = h.bookingService.CreateTicketType(c.Request.Context(), services.CreateTicketInput{
		EventID:  eventID,
		Name:     input.Name,
		Price:    input.Price,
		Quantity: input.Quantity,
	}, capacity)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(201, gin.H{"message": "ticket type created"})
}

func (h *TicketTypeHandler) Update(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid event id"})
		return
	}

	ticketID, err := uuid.Parse(c.Param("ticket_id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid ticket_id"})
		return
	}

	var input UpdateTicketTypeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "invalid input", "details": err.Error()})
		return
	}

	capacity, err := h.eventsProvider.GetEventCapacity(c.Request.Context(), eventID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	err = h.bookingService.UpdateTicketType(c.Request.Context(), ticketID, eventID, services.UpdateTicketTypeInput{
		Name:     input.Name,
		Price:    input.Price,
		Quantity: input.Quantity,
	}, capacity)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TicketTypeHandler) GetByEventID(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid event id"})
		return
	}

	tickets, err := h.bookingService.GetTicketTypesByEventID(c.Request.Context(), eventID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	resp := make([]TicketTypeResponse, len(tickets))
	for i, tt := range tickets {
		resp[i] = TicketTypeResponse{
			ID:        tt.ID,
			EventID:   tt.EventID,
			Name:      tt.Name,
			Price:     tt.Price,
			Quantity:  tt.Quantity,
			Available: tt.Available,
		}
	}

	c.JSON(200, resp)
}

// func (h *TicketTypeHandler) Delete(c *gin.Context) {
// 	ticketID, err := uuid.Parse(c.Param("ticket_id"))
// 	if err != nil {
// 		c.JSON(400, gin.H{"error": "invalid ticket_id"})
// 		return
// 	}

// 	if err := h.bookingService.DeleteTicketType(c.Request.Context(), ticketID); err != nil {
// 		h.handleError(c, err)
// 		return
// 	}

// 	c.Status(http.StatusNoContent)
// }

func (h *TicketTypeHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperr.ErrConflict):
		c.JSON(409, gin.H{"error": "ticket quantity exceeds event capacity"})
	default:
		apperr.Handle(c, err)
	}
}