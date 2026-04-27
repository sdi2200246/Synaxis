package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	apperr "github.com/sdi2200246/synaxis/internal/error"
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
	baseHandler		 *BaseHandler
	ticketsService   *services.TicketTypeService
}

func NewTicketTypeHandler(ts *services.TicketTypeService , bs *BaseHandler) *TicketTypeHandler {
	return &TicketTypeHandler{baseHandler: bs ,ticketsService: ts}
}

func (h *TicketTypeHandler) Create(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid event id"})
		return
	}

	callerID, err := h.baseHandler.getUserIDFromContext(c)
	if err!= nil {
		h.handleError(c , err)
		return
	}


	var input CreateTicketTypeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "invalid input", "details": err.Error()})
		return
	}

	err = h.ticketsService.CreateTicketType(c.Request.Context(),callerID , services.CreateTicketInput{
		EventID:  eventID,
		Name:     input.Name,
		Price:    input.Price,
		Quantity: input.Quantity,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(201, gin.H{"message": "ticket type created"})
}

func (h *TicketTypeHandler) Update(c *gin.Context) {
	callerID, err := h.baseHandler.getUserIDFromContext(c)
	if err!= nil {
		h.handleError(c , err)
		return
	}

	
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

	err = h.ticketsService.UpdateTicketType(c.Request.Context(),callerID , ticketID, eventID, services.UpdateTicketTypeInput{
		Name:     input.Name,
		Price:    input.Price,
		Quantity: input.Quantity,
	})
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

	tickets, err := h.ticketsService.GetTicketTypesByEventID(c.Request.Context(), eventID)
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

func (h *TicketTypeHandler) GetByID(c *gin.Context) {
	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid ticket id"})
		return
	}

	tt, err := h.ticketsService.GetByID(c.Request.Context(), ticketID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(200, TicketTypeResponse{
		ID:        tt.ID,
		EventID:   tt.EventID,
		Name:      tt.Name,
		Price:     tt.Price,
		Quantity:  tt.Quantity,
		Available: tt.Available,
	})
}

func (h *TicketTypeHandler) handleError(c *gin.Context, err error) {
		apperr.Handle(c, err)
}