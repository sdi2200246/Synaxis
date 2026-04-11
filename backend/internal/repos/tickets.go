package repos

import (
	"context"
	"log/slog"

	"github.com/VauntDev/tqla"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sdi2200246/synaxis/internal/entities"
	apperr "github.com/sdi2200246/synaxis/internal/error"
)

type TicketTypeRepo struct {
	db *pgxpool.Pool
}

func NewTicketTypeRepo(db *pgxpool.Pool) *TicketTypeRepo {
	return &TicketTypeRepo{db}
}

func (r *TicketTypeRepo) Create(ctx context.Context, tt entities.TicketType) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO tickettype (id, event_id, name, price, quantity, available, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, tt.ID, tt.EventID, tt.Name, tt.Price, tt.Quantity, tt.Available, tt.CreatedAt)
	if err != nil {
		slog.Error("TicketTypeRepo.Create failed", "error", err)
		return apperr.ErrInternal
	}
	return nil
}

func (r *TicketTypeRepo) GetByID(ctx context.Context, id uuid.UUID) (entities.TicketType, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, event_id, name, price, quantity, available, created_at
		FROM tickettype WHERE id = $1
	`, id)

	var tt entities.TicketType
	err := row.Scan(&tt.ID, &tt.EventID, &tt.Name, &tt.Price, &tt.Quantity, &tt.Available, &tt.CreatedAt)
	if err != nil {
		return entities.TicketType{}, apperr.ErrNotFound
	}
	return tt, nil
}

func (r *TicketTypeRepo) GetByEventID(ctx context.Context, eventID uuid.UUID) ([]entities.TicketType, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, event_id, name, price, quantity, available, created_at
		FROM tickettype WHERE event_id = $1
	`, eventID)
	if err != nil {
		slog.Error("TicketTypeRepo.GetByEventID failed", "error", err)
		return nil, apperr.ErrInternal
	}
	defer rows.Close()

	var result []entities.TicketType
	for rows.Next() {
		var tt entities.TicketType
		if err := rows.Scan(&tt.ID, &tt.EventID, &tt.Name, &tt.Price, &tt.Quantity, &tt.Available, &tt.CreatedAt); err != nil {
			slog.Error("TicketTypeRepo.GetByEventID scan failed", "error", err)
			return nil, apperr.ErrInternal
		}
		result = append(result, tt)
	}
	return result, nil
}

func (r *TicketTypeRepo) SumQuantityByEventID(ctx context.Context, eventID uuid.UUID) (int, error) {
	var sum int
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(quantity), 0) FROM tickettype WHERE event_id = $1
	`, eventID).Scan(&sum)
	if err != nil {
		slog.Error("TicketTypeRepo.SumQuantityByEventID failed", "error", err)
		return 0, apperr.ErrInternal
	}
	return sum, nil
}

func (r *TicketTypeRepo) Update(ctx context.Context, id uuid.UUID, update entities.UpdateTicketType) error {
	t, err := tqla.New(tqla.WithPlaceHolder(tqla.Dollar))
	if err != nil {
		return apperr.ErrInternal
	}

	query, args, err := t.Compile(`
		UPDATE tickettype SET
			{{ if .Update.Name }}     name  = {{ .Update.Name }},  {{ end }}
			{{ if .Update.Price }}    price = {{ .Update.Price }}, {{ end }}
			{{ if .Update.Quantity }} quantity = {{ .Update.Quantity }}, available = available + ({{ .Update.Quantity }} - quantity), {{ end }}
			updated_at = now()
		WHERE id = {{ .ID }}
	`, struct {
		Update entities.UpdateTicketType
		ID     uuid.UUID
	}{update, id})
	if err != nil {
		slog.Error("TicketTypeRepo.Update template failed", "error", err)
		return apperr.ErrInternal
	}

	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		slog.Error("TicketTypeRepo.Update exec failed", "error", err)
		return apperr.ErrInternal
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}

func (r *TicketTypeRepo) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM tickettype WHERE id = $1`, id)
	if err != nil {
		slog.Error("TicketTypeRepo.Delete failed", "error", err)
		return apperr.ErrInternal
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}