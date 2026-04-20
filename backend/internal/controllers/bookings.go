package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/services"
)

type CreateBookingRequest struct {
	TicketTypeID uuid.UUID `json:"ticket_type_id" binding:"required"`
	Quantity     int       `json:"quantity"       binding:"required,min=1"`
}

type BookingHandler struct {
	bookingService *services.BookingService
}

func NewBookingHandler(bs *services.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bs}
}
func (h *BookingHandler) Create(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	userID, ok := val.(uuid.UUID)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid user ID in token"})
		return
	}

	var input CreateBookingRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "invalid input", "details": err.Error()})
		return
	}

	err := h.bookingService.CreateBooking(c.Request.Context(), services.CreateBookingInput{
		TicketTypeID: input.TicketTypeID,
		UserID:       userID,
		Quantity:     input.Quantity,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(201, gin.H{"message": "booking created"})
}

func (h *BookingHandler) GetUserBookings(c *gin.Context) {
	userIDParam, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	bookings, err := h.bookingService.GetUserBookings(c.Request.Context(), userIDParam)
	if err != nil {
		apperr.Handle(c, err)
		return
	}

	c.JSON(200, ToBookingListResponse(bookings))
}

func (h *BookingHandler) GetEventBookings(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid event id"})
		return
	}

	bookings, err := h.bookingService.GetEventBookings(c.Request.Context(), eventID)
	if err != nil {
		apperr.Handle(c, err)
		return
	}

	c.JSON(200, ToBookingListResponse(bookings))
}

func (h *BookingHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperr.ErrConflict):
		c.JSON(409, gin.H{"error": "not enough tickets available"})
	case errors.Is(err, apperr.ErrNotFound):
		c.JSON(404, gin.H{"error": "ticket type not found"})
	default:
		apperr.Handle(c, err)
	}
}