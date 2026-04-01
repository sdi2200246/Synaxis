package repos

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sdi2200246/synaxis/internal/entities"
	"github.com/sdi2200246/synaxis/internal/error"
)

type EventRepo struct{
	db *pgxpool.Pool
}

func NewEventRepo(db *pgxpool.Pool)*EventRepo{
	return  &EventRepo{db}
}

func (r *EventRepo)Create(ctx context.Context , event entities.Event)error{
	_, err := r.db.Exec(ctx, `
		INSERT INTO "event" (
			id, organizer_id, venue_id, title, event_type,
			status, description, capacity,
			start_datetime, end_datetime, created_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8,
			$9, $10, $11
		)`,
		event.ID,
		event.OrganizerID,
		event.VenueID,
		event.Title,
		event.EventType,
		event.Status,
		event.Description,
		event.Capacity,
		event.StartDatetime,
		event.EndDatetime,
		event.CreatedAt,
	)

	if err != nil {
		slog.Error("EventRepo.Create failed", "error", err)
		
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperr.ErrConflict
		}
		return apperr.ErrInternal
	}
	return nil
}


func (r *EventRepo)GetByID(ctx context.Context , id uuid.UUID)(entities.Event , error){
	row := r.db.QueryRow(ctx, `
		SELECT id, organizer_id, venue_id, title, event_type,
			status, description, capacity,
			start_datetime, end_datetime, created_at
		FROM "event"
		WHERE id = $1
	`, id)

	var e entities.Event
	err := row.Scan(
		&e.ID,
		&e.OrganizerID,
		&e.VenueID,
		&e.Title,
		&e.EventType,
		&e.Status,
		&e.Description,
		&e.Capacity,
		&e.StartDatetime,
		&e.EndDatetime,
		&e.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entities.Event{}, apperr.ErrNotFound
		}
		slog.Error("EventRepo.GetByID failed", "error", err)
		return entities.Event{}, apperr.ErrInternal
	}
	return e, nil
}

func (r *EventRepo)GetByOrganizerID(ctx context.Context, organizerID uuid.UUID) ([]entities.Event, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, organizer_id, venue_id, title, event_type,
			status, description, capacity,
			start_datetime, end_datetime, created_at
		FROM "event"
		WHERE organizer_id = $1
	`, organizerID)
	if err != nil {
		slog.Error("EventRepo.GetByOrganizerID failed", "error", err)
		return nil, apperr.ErrInternal
	}
	defer rows.Close()

	var events []entities.Event
	for rows.Next() {
		var e entities.Event
		err := rows.Scan(
			&e.ID,
			&e.OrganizerID,
			&e.VenueID,
			&e.Title,
			&e.EventType,
			&e.Status,
			&e.Description,
			&e.Capacity,
			&e.StartDatetime,
			&e.EndDatetime,
			&e.CreatedAt,
		)
		if err != nil {
			slog.Error("EventRepo.GetByOrganizerID scan failed", "error", err)
			return nil, apperr.ErrInternal
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		slog.Error("EventRepo.GetByOrganizerID rows error", "error", err)
		return nil, apperr.ErrInternal
	}

	return events, nil
}

func (r *EventRepo)Publish(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `
		UPDATE "event"
		SET status = 'PUBLISHED'
		WHERE id = $1 AND status = 'DRAFT'
	`, id)
	if err != nil {
		slog.Error("EventRepo.Publish failed", "error", err)
		return apperr.ErrInternal
	}

	if result.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}

	return nil
}
func (r *EventRepo) Cancel(ctx context.Context, id uuid.UUID) error {
    result, err := r.db.Exec(ctx, `
        UPDATE "event"
        SET status = 'CANCELLED'
        WHERE id = $1 AND status = 'PUBLISHED'
    `, id)
    if err != nil {
        slog.Error("EventRepo.Cancel: update failed", "error", err)
        return apperr.ErrInternal
    }

    if result.RowsAffected() == 0 {
        return apperr.ErrNotFound
    }

    return nil
}