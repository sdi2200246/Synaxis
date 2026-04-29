package repos

import (
	"context"
	"log/slog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sdi2200246/synaxis/internal/entities"
	"github.com/sdi2200246/synaxis/internal/error"
)



type VisitsRepo struct{
	db *pgxpool.Pool
}

func NewVisitsRepo(db *pgxpool.Pool)*VisitsRepo{
    return  &VisitsRepo{db}
}

func (r *VisitsRepo) Create(ctx context.Context, visit entities.Visit) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO visit (id, user_id, event_id, visited_at)
		VALUES ($1, $2, $3, $4)
	`, visit.ID, visit.UserID, visit.EventID, visit.VisitedAt)
	if err != nil {
		slog.Error("VisitsRepo.Create failed", "error", err)
		return apperr.ErrInternal
	}
	return nil
}
