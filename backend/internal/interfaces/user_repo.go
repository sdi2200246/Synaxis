package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/entities"
)

type UserRepository interface {
    Create(ctx context.Context, user entities.User) error
    GetByID(ctx context.Context, id uuid.UUID) (entities.User, error)
    GetByUsername(ctx context.Context, username string) (entities.User, error)
	ListUsers(ctx context.Context , filter entities.UserFilter)([]entities.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, u entities.UserUpdate) error 
}