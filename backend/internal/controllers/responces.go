package controllers

import (
    "time"
    "github.com/google/uuid"
    "github.com/sdi2200246/synaxis/internal/services"
)

type VenueResponse struct {
    ID        uuid.UUID `json:"id"`
    Name      string    `json:"name"`
    Address   string    `json:"address"`
    City      string    `json:"city"`
    Country   string    `json:"country"`
    Latitude  *float64  `json:"latitude,omitempty"`
    Longitude *float64  `json:"longitude,omitempty"`
    Capacity *int       `json:"capacity"`    
}

type CategoryResponce struct{
    ID  uuid.UUID `json:"id"`
    Name string   `json:"name"`
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
    Categories   *[]CategoryResponce `json:"categories"`
}

type AdminUserResponse struct {
    ID        uuid.UUID `json:"id"`
    Username  string    `json:"username"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Email     string    `json:"email"`
    Address   string    `json:"address"`
    City      string    `json:"city"`
    Country   string    `json:"country"`
    TaxID     string    `json:"tax_id"`
    Status    string    `json:"status"`
    Phone     string    `json:"phone"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt *time.Time `json:"updated_at"`
}


func ToEventResponse(ev services.DetailedEvent) EventResponse {

    categories := make([]CategoryResponce , 0)

    for _,c:= range ev.Categories{
        categories = append(categories,  CategoryResponce{ID:c.ID,Name: c.Name})
    }

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
            ID:        ev.VenueID,
            Name:      ev.VenueName,
            Address:   ev.VenueAddress,
            City:      ev.VenueCity,
            Country:   ev.VenueCountry,
            Latitude:  ev.VenueLatitude,
            Longitude: ev.VenueLongitude,
        },
        Categories: &categories,
    }
}

func ToEventListResponse(events []services.DetailedEvent) []EventResponse {
    result := make([]EventResponse, len(events))
    for i, ev := range events {
        result[i] = ToEventResponse(ev)
    }
    return result
}