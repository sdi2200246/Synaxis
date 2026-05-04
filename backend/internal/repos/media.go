package repos

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sdi2200246/synaxis/internal/entities"
	apperr "github.com/sdi2200246/synaxis/internal/error"
)

type MediaRepo struct {
	db *pgxpool.Pool
}

func NewMediaRepo(db *pgxpool.Pool) *MediaRepo {
	return &MediaRepo{db}
}

func (r *MediaRepo) Create(ctx context.Context, media entities.Media) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO media (id, event_id, filename, size_bytes , uploaded_at)
		VALUES ($1, $2, $3, $4 , $5)
	`, media.ID, media.EventID, media.Filename, media.SizeBytes ,media.UploadedAt)
	if err != nil {
		slog.Error("MediaRepo.Create failed", "error", err)
		return apperr.ErrInternal
	}
	return nil
}

func (r *MediaRepo) GetByEventID(ctx context.Context, eventID uuid.UUID) ([]entities.Media, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, event_id, filename, uploaded_at
		FROM media WHERE event_id = $1
	`, eventID)
	if err != nil {
		slog.Error("MediaRepo.GetByEventID failed", "error", err)
		return nil, apperr.ErrInternal
	}
	defer rows.Close()

	var result []entities.Media
	for rows.Next() {
		var m entities.Media
		if err := rows.Scan(&m.ID, &m.EventID, &m.Filename, &m.UploadedAt); err != nil {
			slog.Error("MediaRepo.GetByEventID scan failed", "error", err)
			return nil, apperr.ErrInternal
		}
		result = append(result, m)
	}
	if err := rows.Err(); err != nil {
		return nil, apperr.ErrInternal
	}
	return result, nil
}

func (r *MediaRepo) GetByEventIDs(ctx context.Context, eventIDs []uuid.UUID) (map[uuid.UUID][]entities.Media, error) {
	if len(eventIDs) == 0 {
		return map[uuid.UUID][]entities.Media{}, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, event_id, filename, uploaded_at
		FROM media
		WHERE event_id = ANY($1)
	`, eventIDs)
	if err != nil {
		slog.Error("MediaRepo.GetByEventIDs failed", "error", err)
		return nil, apperr.ErrInternal
	}
	defer rows.Close()

	result := make(map[uuid.UUID][]entities.Media)
	for rows.Next() {
		var m entities.Media
		if err := rows.Scan(&m.ID, &m.EventID, &m.Filename, &m.UploadedAt); err != nil {
			slog.Error("MediaRepo.GetByEventIDs scan failed", "error", err)
			return nil, apperr.ErrInternal
		}
		result[m.EventID] = append(result[m.EventID], m)
	}
	if err := rows.Err(); err != nil {
		return nil, apperr.ErrInternal
	}
	return result, nil
}

func (r *MediaRepo) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM media WHERE id = $1`, id)
	if err != nil {
		slog.Error("MediaRepo.Delete failed", "error", err)
		return apperr.ErrInternal
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}