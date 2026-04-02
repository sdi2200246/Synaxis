package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/services"
)
type EventsHandler struct {
    eventsService *services.EventService
}

func NewEventsHandler(eventsService *services.EventService) *EventsHandler {
    return &EventsHandler{eventsService: eventsService}
}

func (h *EventsHandler)Create(c *gin.Context) {

    val, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	organizerID, ok := val.(uuid.UUID)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid user ID in token"})
		return
	}

	var input services.CandidateEvent
    if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "invalid input", "details": err.Error()})
        return
    }

    err := h.eventsService.CreateEvent(c.Request.Context(), organizerID , input)
    if err != nil {
		h.handleError(c , err)
        return
    }
    c.JSON(201, gin.H{"message": "draft event created succesfully"})
}


func (h *EventsHandler) GetMyEvents(c *gin.Context) {
   
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	
	organizerID, ok := val.(uuid.UUID)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid user ID in token"})
		return
	}

    events, err := h.eventsService.GetOrganizerEvents(c.Request.Context(), organizerID)
    if err != nil {
        apperr.Handle(c, err)
        return
    }

    c.JSON(200, ToEventListResponse(events))
}


func (h *EventsHandler) handleError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, apperr.ErrConflict):
        c.JSON(409, gin.H{"error": err.Error(), "fields": "starting_time"})
    default:
        apperr.Handle(c, err)
    }
}