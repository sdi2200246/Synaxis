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

type VenueRepo struct{
	db *pgxpool.Pool
}

func NewVenueRepo(db *pgxpool.Pool)*VenueRepo{
    return  &VenueRepo{db}
}

func (r *VenueRepo) ListVenues(ctx context.Context, filter entities.VenuesFilter) ([]entities.Venue, error) {
    t, err := tqla.New(tqla.WithPlaceHolder(tqla.Dollar))
    if err != nil {
        return nil, apperr.ErrInternal
    }
	
    query, args, err := t.Compile(`
        SELECT 
            id, name, address, city, country, latitude, longitude, capacity
        FROM venue
        WHERE 1=1
        {{ if .Name }} AND name ILIKE '%' || {{ .Name }} || '%' {{ end }}
        {{ if .Capacity }} AND capacity >= {{ .Capacity }} {{ end }}
    `, filter)
    
    if err != nil {
        slog.Error("ListVenues template failed", "error", err)
        return nil, apperr.ErrInternal
    }

    rows, err := r.db.Query(ctx, query, args...)
    if err != nil {
        slog.Error("ListVenues query failed", "error", err)
        return nil, apperr.ErrInternal
    }
    defer rows.Close()

    var venues []entities.Venue
    for rows.Next() {
        var v entities.Venue
        err := rows.Scan(
            &v.ID,
            &v.Name,
            &v.Address,
            &v.City,
            &v.Country,
            &v.Latitude,
            &v.Longitude,
            &v.Capacity,
        )
        if err != nil {
            slog.Error("ListVenues scan failed", "error", err)
            return nil, apperr.ErrInternal
        }
        venues = append(venues, v)
    }

    if err := rows.Err(); err != nil {
        slog.Error("ListVenues iteration failed", "error", err)
        return nil, apperr.ErrInternal
    }

    return venues, nil
}
func (r *VenueRepo) GetByID(ctx context.Context, id uuid.UUID) (entities.Venue, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, name, address, city, country, latitude, longitude, capacity
		FROM venue
		WHERE id = $1
	`, id)

	var v entities.Venue
	err := row.Scan(
		&v.ID,
		&v.Name,
		&v.Address,
		&v.City,
		&v.Country,
		&v.Latitude,
		&v.Longitude,
		&v.Capacity,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entities.Venue{}, apperr.ErrNotFound
		}
		slog.Error("VenueRepo.GetByID failed", "error", err)
		return entities.Venue{}, apperr.ErrInternal
	}
	return v, nil
}