package interfaces

import (
	"context"
	"github.com/google/uuid"
)

type EventCapacityProvider interface {
	GetEventCapacity(ctx context.Context, eventID uuid.UUID) (int, error)
}