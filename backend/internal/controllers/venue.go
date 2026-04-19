package controllers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/services"
)

type VenueFilter struct {
	Name     string `form:"name"`
	Capacity string `form:"capacity"`
}

type VenueHandler struct {
	venueService *services.VenueService
}

func NewVenueHandler(venueService *services.VenueService) *VenueHandler {
	return &VenueHandler{venueService: venueService}
}

func (h *VenueHandler) GetVenues(c *gin.Context) {
	var filter VenueFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(400, gin.H{"error": "invalid query params"})
		return
	}

	svcFilter := services.VenueFilter{}

	if filter.Name != "" {
		svcFilter.Name = &filter.Name
	}

	if filter.Capacity != "" {
		capVal, err := strconv.Atoi(filter.Capacity)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid capacity format"})
			return
		}
		svcFilter.Capacity = &capVal
	}

	venues, err := h.venueService.GetVenues(c.Request.Context(), svcFilter)
	if err != nil {
		h.handleError(c, err)
		return
	}

	plain := make([]VenueResponse, len(venues))
	for i, v := range venues {
		plain[i] = VenueResponse{
			ID:       v.ID,
			Name:     v.Name,
			City:     v.City,
			Country:  v.Country,
			Capacity: v.Capacity,
		}
	}
	c.JSON(200, gin.H{"count": len(venues), "venues": plain})
}


func (h *VenueHandler) GetVenue(c *gin.Context) {
	venueID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid venue id"})
		return
	}

	venue, err := h.venueService.GetVenue(c.Request.Context(), venueID)
	if err != nil {
		apperr.Handle(c, err)
		return
	}

	c.JSON(200, ToVenueResponse(venue))
}

func (h *VenueHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperr.ErrBadInput):
		c.JSON(400, gin.H{"error": err.Error()})
	case errors.Is(err, apperr.ErrNotFound):
		c.JSON(404, gin.H{"error": "venue not found"})
	default:
		apperr.Handle(c, err)
	}
}