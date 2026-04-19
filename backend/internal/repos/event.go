package repos

import (
	"context"
	"log/slog"

	"github.com/VauntDev/tqla"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (r *EventRepo) CreateWithCategories(ctx context.Context, event entities.Event, categoryIDs []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
        return apperr.ErrInternal
    }
    defer tx.Rollback(ctx)

    _, err = tx.Exec(ctx, `
        INSERT INTO event (
            id, organizer_id, venue_id, title, event_type,
            status, description, capacity, start_datetime, end_datetime, created_at
        ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
    `,
        event.ID, event.OrganizerID, event.VenueID, event.Title, event.EventType,
        event.Status, event.Description, event.Capacity, event.StartDatetime,
        event.EndDatetime, event.CreatedAt,
    )
    if err != nil {
        return apperr.ErrInternal
    }

    for _, catID := range categoryIDs {
        _, err = tx.Exec(ctx, `
            INSERT INTO eventcategory (event_id, category_id)
            VALUES ($1, $2)
        `, event.ID, catID)
        if err != nil {
            return apperr.ErrInternal
        }
    }

    return tx.Commit(ctx)
}

func (r *EventRepo) Delete(ctx context.Context, eventID uuid.UUID) error {
    cmd, err := r.db.Exec(ctx, `DELETE FROM event WHERE id = $1`, eventID)
    if err != nil {
        return apperr.ErrInternal
    }

    if cmd.RowsAffected() == 0 {
        return apperr.ErrNotFound
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


func (r *EventRepo) Update(ctx context.Context, eventID uuid.UUID, update entities.UpdateEvent) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return apperr.ErrInternal
	}
	defer tx.Rollback(ctx)

	t, err := tqla.New(tqla.WithPlaceHolder(tqla.Dollar))
	if err != nil {
		return apperr.ErrInternal
	}

	query, args, err := t.Compile(`
		UPDATE event SET
			{{ if .Title }} title = {{ .Title }}, {{ end }}
			{{ if .EventType }} event_type = {{ .EventType }}, {{ end }}
			{{ if .VenueID }} venue_id = {{ .VenueID }}, {{ end }}
			{{ if .Description }} description = {{ .Description }}, {{ end }}
			{{ if .Status }} status = {{ .Status }}, {{ end }}
			updated_at = now()
		WHERE id = {{ .EventID }}
	`, struct {
		entities.UpdateEvent
		EventID uuid.UUID
	}{update, eventID})

	if err != nil {
		slog.Error("Update event template failed", "error", err)
		return apperr.ErrInternal
	}

	tag, err := tx.Exec(ctx, query, args...)
	if err != nil {
		slog.Error("Update event query failed", "error", err)
		return apperr.ErrInternal
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}

	if update.CategoryIDs != nil {
		_, err = tx.Exec(ctx, `DELETE FROM eventcategory WHERE event_id = $1`, eventID)
		if err != nil {
			slog.Error("Update event delete categories failed", "error", err)
			return apperr.ErrInternal
		}
		for _, catID := range *update.CategoryIDs {
			_, err = tx.Exec(ctx, `
				INSERT INTO eventcategory (event_id, category_id) VALUES ($1, $2)
			`, eventID, catID)
			if err != nil {
				slog.Error("Update event insert category failed", "error", err)
				return apperr.ErrInternal
			}
		}
	}
    if err := tx.Commit(ctx); err != nil {
        slog.Error("Update event commit failed", "error", err)
        return apperr.ErrInternal
    }
    
    return nil

}

func (r *EventRepo) GetbyFilter(ctx context.Context, filter entities.EventFilter) ([]entities.Event, bool, error) {
    if filter.Limit == 0 {
        filter.Limit = 20
    }
    if filter.Limit > 100 {
        filter.Limit = 100
    }

    t, err := tqla.New(tqla.WithPlaceHolder(tqla.Dollar))
    if err != nil {
        return nil, false, apperr.ErrInternal
    }

    query, args, err := t.Compile(`
		SELECT e.id, e.organizer_id, e.venue_id, e.title, e.event_type,
			e.status, e.description, e.capacity,
			e.start_datetime, e.end_datetime, e.created_at
		FROM event e
		JOIN venue v ON e.venue_id = v.id
		WHERE 1=1
		{{ if .Status }} AND e.status = {{ .Status }} {{ end }}
		{{ if .OrganizerID }} AND e.organizer_id = {{ .OrganizerID }} {{ end }}
		{{ if or .Title .Description }}
		AND (
			1=0
			{{ if .Title }} OR e.title ILIKE '%' || {{ .Title }} || '%' {{ end }}
			{{ if .Description }} OR e.description ILIKE '%' || {{ .Description }} || '%' {{ end }}
		)
		{{ end }}
		{{ if .City }} AND v.city ILIKE '%' || {{ .City }} || '%' {{ end }}
		{{ if .Country }} AND v.country ILIKE '%' || {{ .Country }} || '%' {{ end }}
		{{ if .StartAfter }} AND e.start_datetime >= {{ .StartAfter }} {{ end }}
		{{ if .StartBefore }} AND e.start_datetime <= {{ .StartBefore }} {{ end }}
		{{ if .CategoryIDs }} AND e.id IN (
			SELECT event_id FROM eventcategory WHERE category_id = ANY({{ .CategoryIDs }})
		) {{ end }}
		{{ if .MinPrice }} AND e.id IN (
			SELECT event_id FROM tickettype WHERE price >= {{ .MinPrice }}
		) {{ end }}
		{{ if .MaxPrice }} AND e.id IN (
			SELECT event_id FROM tickettype WHERE price <= {{ .MaxPrice }}
		) {{ end }}
		ORDER BY e.start_datetime ASC
		LIMIT {{ .Limit }} OFFSET {{ .Offset }}
	`, filter)
    if err != nil {
        slog.Error("GetbyFilter template failed", "error", err)
        return nil, false, apperr.ErrInternal
    }

    rows, err := r.db.Query(ctx, query, args...)
    if err != nil {
        slog.Error("GetbyFilter query failed", "error", err)
        return nil, false, apperr.ErrInternal
    }
    defer rows.Close()

    var results []entities.Event
    for rows.Next() {
        var e entities.Event
        if err := rows.Scan(
            &e.ID, &e.OrganizerID, &e.VenueID, &e.Title, &e.EventType,
            &e.Status, &e.Description, &e.Capacity,
            &e.StartDatetime, &e.EndDatetime, &e.CreatedAt,
        ); err != nil {
            slog.Error("GetbyFilter scan failed", "error", err)
            return nil, false, apperr.ErrInternal
        }
        results = append(results, e)
    }

    if err := rows.Err(); err != nil {
        slog.Error("SearchPublished iteration failed", "error", err)
        return nil, false, apperr.ErrInternal
    }

    hasMore := len(results) == filter.Limit
    return results, hasMore, nil
}


func (r *EventRepo) GetAll(ctx context.Context) ([]entities.Event, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, organizer_id, venue_id, title, event_type,
			status, description, capacity,
			start_datetime, end_datetime, created_at
		FROM event
	`)
	if err != nil {
		slog.Error("EventRepo.GetAll failed", "error", err)
		return nil, apperr.ErrInternal
	}
	defer rows.Close()

	var results []entities.Event
	for rows.Next() {
		var e entities.Event
		if err := rows.Scan(
			&e.ID, &e.OrganizerID, &e.VenueID, &e.Title, &e.EventType,
			&e.Status, &e.Description, &e.Capacity,
			&e.StartDatetime, &e.EndDatetime, &e.CreatedAt,
		); err != nil {
			slog.Error("EventRepo.GetAll scan failed", "error", err)
			return nil, apperr.ErrInternal
		}
		results = append(results, e)
	}

	if err := rows.Err(); err != nil {
		slog.Error("EventRepo.GetAll rows error", "error", err)
		return nil, apperr.ErrInternal
	}

	return results, nil
}