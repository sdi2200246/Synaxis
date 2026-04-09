package interfaces

import (
    "context"
    "github.com/sdi2200246/synaxis/internal/entities"
    "github.com/google/uuid"
)

type EventRepository interface {
    CreateWithCategories(ctx context.Context, event entities.Event ,categoryIDs []uuid.UUID) error
    GetByID(ctx context.Context, id uuid.UUID) (entities.Event, error)
    GetByOrganizerID(ctx context.Context, organizerID uuid.UUID) ([]entities.OrganizerEvent, error) 
    Update(ctx context.Context, eventID uuid.UUID, update entities.UpdateEvent) error
	// Publish(ctx context.Context, id uuid.UUID) error 
	// Cancel(ctx context.Context, id uuid.UUID) error
}