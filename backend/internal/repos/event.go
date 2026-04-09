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


func (r *EventRepo) GetByOrganizerID(ctx context.Context, organizerID uuid.UUID) ([]entities.OrganizerEvent, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			e.id, e.organizer_id, e.venue_id, e.title, e.event_type,
			e.status, e.description, e.capacity,
			e.start_datetime, e.end_datetime, e.created_at,
			v.id, v.name, v.address, v.city, v.country,
			v.latitude, v.longitude, v.capacity,
			array_agg(c.id) as category_ids,
			array_agg(c.name) as category_names
		FROM event e
		JOIN venue v ON e.venue_id = v.id
		JOIN eventcategory ec ON ec.event_id = e.id
		JOIN category c ON c.id = ec.category_id
		WHERE e.organizer_id = $1
		GROUP BY e.id, v.id
	`, organizerID)
	if err != nil {
		slog.Error("EventRepo.GetByOrganizerID failed", "error", err)
		return nil, apperr.ErrInternal
	}
	defer rows.Close()

	var results []entities.OrganizerEvent
	for rows.Next() {
		var ev entities.OrganizerEvent
		var categoryIDs []uuid.UUID
		var categoryNames []string

		err := rows.Scan(
			&ev.ID, &ev.OrganizerID, &ev.VenueID, &ev.Title, &ev.EventType,
			&ev.Status, &ev.Description, &ev.Capacity,
			&ev.StartDatetime, &ev.EndDatetime, &ev.CreatedAt,
			&ev.Venue.ID, &ev.Venue.Name, &ev.Venue.Address, &ev.Venue.City,
			&ev.Venue.Country, &ev.Venue.Latitude, &ev.Venue.Longitude, &ev.Venue.Capacity,
			&categoryIDs, &categoryNames,
		)
		if err != nil {
			slog.Error("EventRepo.GetByOrganizerID scan failed", "error", err)
			return nil, apperr.ErrInternal
		}

		for i := range categoryIDs {
			ev.Categories = append(ev.Categories, entities.Category{
				ID:   categoryIDs[i],
				Name: categoryNames[i],
			})
		}

		results = append(results, ev)
	}

	if err := rows.Err(); err != nil {
		slog.Error("EventRepo.GetByOrganizerID rows error", "error", err)
		return nil, apperr.ErrInternal
	}
	
	return results, nil
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