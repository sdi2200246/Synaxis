package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/entities"
	"github.com/sdi2200246/synaxis/internal/interfaces"
	"github.com/sdi2200246/synaxis/internal/repos"
)

type VenueFilter struct{
	Name *string 
	Capacity *int
}

type Venue struct{
	ID       uuid.UUID 
    Name     string
	City     string
	Country  string
    Capacity  *int
} 

type DetailedVenue struct {
	ID        uuid.UUID
	Name      string
	Address   string
	City      string
	Country   string
	Latitude  *float64
	Longitude *float64
	Capacity  *int
}


type VenueService struct{
	venueRepo interfaces.VenuesRepository
}

func NewVenueService(db *repos.VenueRepo)*VenueService{
	return &VenueService{db}

}

func (s *VenueService) GetVenues(ctx context.Context, f VenueFilter) ([]Venue, error) {
    filter := entities.VenuesFilter{
        Name:     f.Name,
        Capacity: f.Capacity,
    }

    dbVenues, err := s.venueRepo.ListVenues(ctx, filter)
    if err != nil {
        return nil, err
    }

    venues := make([]Venue, len(dbVenues))
    for i, v := range dbVenues {
        venues[i] = Venue{
            ID:       v.ID,
            Name:     v.Name,
			City:     v.City,
			Country:  v.Country,
            Capacity: v.Capacity, 
        }
    }

    return venues, nil
}

func (s *VenueService) GetVenue(ctx context.Context, id uuid.UUID) (DetailedVenue, error) {
	v, err := s.venueRepo.GetByID(ctx, id)
	if err != nil {
		return DetailedVenue{}, err
	}

	return DetailedVenue{
		ID:        v.ID,
		Name:      v.Name,
		Address:   v.Address,
		City:      v.City,
		Country:   v.Country,
		Latitude:  v.Latitude,
		Longitude: v.Longitude,
		Capacity:  v.Capacity,
	}, nil
}