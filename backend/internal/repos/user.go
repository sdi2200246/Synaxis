package repos

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sdi2200246/synaxis/internal/entities"
	"github.com/sdi2200246/synaxis/internal/error"
    "github.com/google/uuid"
)

type UserRepo struct{
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool)*UserRepo{
    return  &UserRepo{db}

}
func (r *UserRepo)Create(ctx context.Context, user entities.User) error {
    _, err := r.db.Exec(ctx, `
        INSERT INTO "user" (
            id, username, password_hash, first_name, last_name,
            email, phone, address, city, country, tax_id,
            role, status, created_at
        ) VALUES (
            $1, $2, $3, $4, $5,
            $6, $7, $8, $9, $10, $11,
            $12, $13, $14
        )`,
        user.ID,
        user.Username,
        user.PasswordHash,
        user.FirstName,
        user.LastName,
        user.Email,
        user.Phone,
        user.Address,
        user.City,
        user.Country,
        user.TaxID,
        user.Role,
        user.Status,
        user.CreatedAt,
    )
    
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            if strings.Contains(pgErr.ConstraintName, "username") {
                return apperr.ErrUsernameConflict
            }
            if strings.Contains(pgErr.ConstraintName, "email") {
                return apperr.ErrEmailConflict
            }
            return apperr.ErrConflict
        }
        slog.Error("UserRepository.Create failed", "error", err)
        return apperr.ErrInternal
    }
    return nil
}

func (r *UserRepo)GetByID(ctx context.Context, id uuid.UUID) (entities.User, error) {    
    row := r.db.QueryRow(ctx, `
        SELECT id, username, password_hash, first_name, last_name,
            email, phone, address, city, country, tax_id,
            role, status, created_at
        FROM "user"
        WHERE id = $1
    `, id)

    var u entities.User
    err := row.Scan(
        &u.ID,
        &u.Username,
        &u.PasswordHash,
        &u.FirstName,
        &u.LastName,
        &u.Email,
        &u.Phone,
        &u.Address,
        &u.City,
        &u.Country,
        &u.TaxID,
        &u.Role,
        &u.Status,
        &u.CreatedAt,
    )
    if err != nil {
        if err == pgx.ErrNoRows {
            return entities.User{}, apperr.ErrNotFound
        }
        slog.Error("UserRepo.FindByID failed", "error", err)
        return entities.User{}, apperr.ErrInternal
    }
    return u, nil
}
func (r *UserRepo)GetByUsername(ctx context.Context, username string) (entities.User, error) {
    row := r.db.QueryRow(ctx, `
        SELECT id, username, password_hash, first_name, last_name,
               email, phone, address, city, country, tax_id,
               role, status, created_at
        FROM "user"
        WHERE username = $1
    `, username)

    var u entities.User
    err := row.Scan(
        &u.ID,
        &u.Username,
        &u.PasswordHash,
        &u.FirstName,
        &u.LastName,
        &u.Email,
        &u.Phone,
        &u.Address,
        &u.City,
        &u.Country,
        &u.TaxID,
        &u.Role,
        &u.Status,
        &u.CreatedAt,
    )
    if err != nil {
        if err == pgx.ErrNoRows {
            return entities.User{}, apperr.ErrNotFound
        }
        slog.Error("UserRepo.FindByUsername failed", "error", err)
        return entities.User{}, apperr.ErrInternal
    }
    return u, nil
}