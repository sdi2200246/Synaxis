package interfaces

import (
    "context"
    "github.com/sdi2200246/synaxis/internal/entities"
    "github.com/google/uuid"
)

type EventRepository interface {
    Create(ctx context.Context, event entities.Event) error
    GetByID(ctx context.Context, id uuid.UUID) (entities.Event, error)
	GetByOrganizerID(ctx context.Context, organizerID uuid.UUID) ([]entities.Event, error) 
	Publish(ctx context.Context, id uuid.UUID) error 
	Cancel(ctx context.Context, id uuid.UUID) error
}