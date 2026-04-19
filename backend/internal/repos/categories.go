package repos

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sdi2200246/synaxis/internal/entities"
	apperr "github.com/sdi2200246/synaxis/internal/error"
)

type CategoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) GetAll(ctx context.Context) ([]entities.Category, error) {
	const query = `SELECT id, name FROM category ORDER BY name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		slog.Error("Failed to load categories")
		return nil, apperr.ErrInternal
	}
	defer rows.Close()

	var categories []entities.Category
	for rows.Next() {
		var c entities.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, apperr.ErrInternal
		}
		categories = append(categories, c)
	}

	return categories, rows.Err()
}

func (r *CategoryRepo) GetByEventID(ctx context.Context, eventID uuid.UUID) ([]entities.Category, error) {
	rows, err := r.db.Query(ctx,
		`SELECT c.id, c.name, c.parent_id
		 FROM category c
		 JOIN eventcategory ec ON ec.category_id = c.id
		 WHERE ec.event_id = $1`,
		eventID,
	)
	if err != nil {
		slog.Error("CategoryRepo.GetByEventID failed", "error", err)
		return nil, apperr.ErrInternal
	}
	defer rows.Close()

	var categories []entities.Category
	for rows.Next() {
		var c entities.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.ParentID); err != nil {
			slog.Error("CategoryRepo.GetByEventID scan failed", "error", err)
			return nil, apperr.ErrInternal
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return nil, apperr.ErrInternal
	}
	return categories, nil
}