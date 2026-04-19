package controllers

import (
	"errors"
	"log/slog"
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/services"
)


type CreateEventRequest struct {
    Title       string    `json:"title"        binding:"required"`
    EventType   string    `json:"event_type"   binding:"required"`
    VenueID     uuid.UUID `json:"venue_id"     binding:"required"`
    Description string    `json:"description"  binding:"required"`
    Capacity    int       `json:"capacity"     binding:"required,min=1"`
    StartDatetime time.Time `json:"start_datetime" binding:"required"`
    EndDatetime   time.Time `json:"end_datetime"   binding:"required"`
    CategoryIDs []uuid.UUID `json:"category_ids" binding:"required,min=1"`
}

type UpdateEventRequest struct{
    EventType   *string      `json:"event_type,omitempty"`
    VenueID     *uuid.UUID   `json:"venue_id,omitempty"`
    Description *string      `json:"description,omitempty"`
    CategoryIDs *[]uuid.UUID `json:"category_ids,omitempty"`
    Status      *string      `json:"status,omitempty"`
}

type SearchEventRequest struct {
    CategoryIDs []string `form:"category_id"`
    Title       *string     `form:"title"`
    Description *string     `form:"description"`
    City        *string     `form:"city"`
    Country     *string     `form:"country"`
    StartAfter  *time.Time  `form:"start_after" time_format:"2006-01-02T15:04:05Z07:00"`
    StartBefore *time.Time  `form:"start_before" time_format:"2006-01-02T15:04:05Z07:00"`
    MinPrice    *float64    `form:"min_price"`
    MaxPrice    *float64    `form:"max_price"`
    Limit       int         `form:"limit,default=20"`
    Offset      int         `form:"offset,default=0"`
}

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

	var input CreateEventRequest
    if err := c.ShouldBindJSON(&input); err != nil {
		slog.Error("Invalid input", "error", err)
		c.JSON(400, gin.H{"error": "invalid input", "details": err.Error()})
        return
    }

    err := h.eventsService.CreateEvent(c.Request.Context(), organizerID , services.CreateEventInput{
		Title: input.Title,
		EventType: input.EventType,
		VenueID: input.VenueID,
		Description: input.Description,
		Capacity: input.Capacity,
		StartDatetime: input.StartDatetime,
		EndDatetime: input.EndDatetime,
		CategoryIDs: input.CategoryIDs,
	})
    if err != nil {
		h.handleError(c , err)
        return
    }
    c.JSON(201, gin.H{"message": "draft event created succesfully"})
}
func (h *EventsHandler)UpdateEvent(c *gin.Context) {

   	eventID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(400, gin.H{"error": "invalid id"})
        return
    }
	
	var input UpdateEventRequest
    if err := c.ShouldBindJSON(&input); err != nil {
	    slog.Error("Invalid input", "error", err)
		c.JSON(400, gin.H{"error": "invalid input", "details": err.Error()})
        return
    }

    err = h.eventsService.UpdateEvent(c.Request.Context(), eventID , services.UpdateEventInput{
		EventType: input.EventType,
		VenueID: input.VenueID,
		Description: input.Description,
		CategoryIDs: input.CategoryIDs,
	})
    if err != nil {
		h.handleError(c , err)
        return
    }
    c.Status(http.StatusNoContent)
}

func (h *EventsHandler) GetOrganizerEvents(c *gin.Context) {
   
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


func (h *EventsHandler) SearchPublished(c *gin.Context) {
    var req SearchEventRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    var categoryIDs []uuid.UUID
    for _, s := range req.CategoryIDs {
        id, err := uuid.Parse(s)
        if err != nil {
            c.JSON(400, gin.H{"error": "invalid category_id: " + s})
            return
        }
        categoryIDs = append(categoryIDs, id)
    }

    filter := services.EventFilterInput{
        CategoryIDs: categoryIDs,
        Title:       req.Title,
        Description: req.Description,
        City:        req.City,
        Country:     req.Country,
        StartAfter:  req.StartAfter,
        StartBefore: req.StartBefore,
        MinPrice:    req.MinPrice,
        MaxPrice:    req.MaxPrice,
        Limit:       req.Limit,
        Offset:      req.Offset,
    }

    events, hasMore, err := h.eventsService.SearchEvents(c.Request.Context(), filter)
    if err != nil {
        apperr.Handle(c, err)
        return
    }

    c.JSON(200, gin.H{
        "events":   ToEventListResponse(events),
        "has_more": hasMore,
    })
}


func (h *EventsHandler) Delete(c *gin.Context) {
    eventID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(400, gin.H{"error": "invalid event id"})
        return
    }

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

    organizerID, err := h.eventsService.GetEventOrganizer(c.Request.Context(), eventID)
    if err != nil {
        apperr.Handle(c, err)
        return
    }
    if organizerID != userID {
        c.JSON(403, gin.H{"error": "forbidden"})
        return
    }

    if err := h.eventsService.Delete(c.Request.Context(), eventID); err != nil {
         c.JSON(403, gin.H{"error": "This event cannot be deleted"})
        return
    }

    c.Status(204)
}