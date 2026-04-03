package controllers

import (
    "time"
    "github.com/google/uuid"
    "github.com/sdi2200246/synaxis/internal/entities"
)

type VenueResponse struct {
    ID        uuid.UUID `json:"id"`
    Name      string    `json:"name"`
    Address   string    `json:"address"`
    City      string    `json:"city"`
    Country   string    `json:"country"`
    Latitude  *float64  `json:"latitude,omitempty"`
    Longitude *float64  `json:"longitude,omitempty"`
}

type EventResponse struct {
    ID            uuid.UUID     `json:"id"`
    Title         string        `json:"title"`
    EventType     string        `json:"event_type"`
    Status        string        `json:"status"`
    Description   string        `json:"description"`
    Capacity      int           `json:"capacity"`
    StartDatetime time.Time     `json:"start_datetime"`
    EndDatetime   time.Time     `json:"end_datetime"`
    Venue         VenueResponse `json:"venue"`
}

func ToEventResponse(ev entities.EventWithVenue) EventResponse {
    return EventResponse{
        ID:            ev.ID,
        Title:         ev.Title,
        EventType:     ev.EventType,
        Status:        ev.Status,
        Description:   ev.Description,
        Capacity:      ev.Capacity,
        StartDatetime: ev.StartDatetime,
        EndDatetime:   ev.EndDatetime,
        Venue: VenueResponse{
            ID:        ev.Venue.ID,
            Name:      ev.Venue.Name,
            Address:   ev.Venue.Address,
            City:      ev.Venue.City,
            Country:   ev.Venue.Country,
            Latitude:  ev.Venue.Latitude,
            Longitude: ev.Venue.Longitude,
        },
    }
}

func ToEventListResponse(events []entities.EventWithVenue) []EventResponse {
    result := make([]EventResponse, len(events))
    for i, ev := range events {
        result[i] = ToEventResponse(ev)
    }
    return result
}