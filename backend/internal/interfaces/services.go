package interfaces

import (
	"context"
	"github.com/google/uuid"
)

type EventsProvider interface {
	GetEventCapacity(ctx context.Context, eventID uuid.UUID) (int, error)
	GetEventStatus(ctx context.Context , id uuid.UUID)(string , error)
	GetEventOrganizer(ctx context.Context , id uuid.UUID)(uuid.UUID , error)
}
