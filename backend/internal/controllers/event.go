package controllers

import (
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
    Title       *string      `json:"title,omitempty"`
    EventType   *string      `json:"event_type,omitempty"`
    VenueID     *uuid.UUID   `json:"venue_id,omitempty"`
    Description *string      `json:"description,omitempty"`
    CategoryIDs *[]uuid.UUID `json:"category_ids,omitempty"`
    Status      *string      `json:"status,omitempty"`
}

type SearchEventRequest struct {
    OrganizerID *string     `form:"organizer_id"`
    Status      *string     `form:"status"`
    CategoryIDs []string    `form:"category_id"`
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
    baseHandler *BaseHandler
    eventsService *services.EventService
}

func NewEventsHandler(eventsService *services.EventService , bh *BaseHandler) *EventsHandler {
    return &EventsHandler{baseHandler: bh , eventsService: eventsService}
}

func (h *EventsHandler)Create(c *gin.Context) {

	userID, err := h.baseHandler.getUserIDFromContext(c)
	if err!= nil {
		h.handleError(c , err)
		return
	}

	var input CreateEventRequest
    if err := c.ShouldBindJSON(&input); err != nil {
		slog.Error("Invalid input", "error", err)
		c.JSON(400, gin.H{"error": "invalid input", "details": err.Error()})
        return
    }

    err = h.eventsService.CreateEvent(c.Request.Context(), userID , services.CreateEventInput{
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

    callerID, err := h.baseHandler.getUserIDFromContext(c)
	if err!= nil {
		h.handleError(c , err)
		return
	}

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

    err = h.eventsService.UpdateEvent(c.Request.Context(), callerID , eventID , services.UpdateEventInput{
        Title:     input.Title,
		EventType: input.EventType,
		VenueID: input.VenueID,
		Description: input.Description,
		CategoryIDs: input.CategoryIDs,
        Status: input.Status,
	})
    if err != nil {
		h.handleError(c , err)
        return
    }
    c.Status(http.StatusNoContent)
}

func (h *EventsHandler) List(c *gin.Context) {
    var req SearchEventRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
   
    var organizerID *uuid.UUID
    if req.OrganizerID != nil {
        id, err := uuid.Parse(*req.OrganizerID)
        if err != nil {
            c.JSON(400, gin.H{"error": "invalid organizer_id"})
            return
        }
        organizerID = &id
    }

    callerID , _ := h.baseHandler.CallerIDExists(c)

    categoryIDs := make([]uuid.UUID, 0, len(req.CategoryIDs))
    for _, s := range req.CategoryIDs {
        id, err := uuid.Parse(s)
        if err != nil {
            c.JSON(400, gin.H{"error": "invalid category_id: " + s})
            return
        }
        categoryIDs = append(categoryIDs, id)
    }

    filter := services.EventFilterInput{
        OrganizerID: organizerID,
        Status:      req.Status,
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

    events, hasMore, err := h.eventsService.List(c.Request.Context(), callerID , filter)
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

    callerID, err := h.baseHandler.getUserIDFromContext(c)
	if err!= nil {
		h.handleError(c , err)
		return
	}
   
    if err := h.eventsService.Delete(c.Request.Context(),callerID ,eventID); err != nil {
        apperr.Handle(c , err)
        return
    }

    c.Status(204)
}

func (h *EventsHandler) GetEventCategories(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid event id"})
		return
	}

	categories, err := h.eventsService.GetEventCategories(c.Request.Context(), eventID)
	if err != nil {
		apperr.Handle(c, err)
		return
	}

	c.JSON(200, ToCategoryListResponse(categories))
}

func (h *EventsHandler) GetByID(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid event id"})
		return
	}

	event, err := h.eventsService.GetByID(c.Request.Context(), eventID)
	if err != nil {
		apperr.Handle(c, err)
		return
	}

	c.JSON(200, ToEventResponse(event))
}

func (h *EventsHandler) handleError(c *gin.Context, err error) {
    apperr.Handle(c, err)
}