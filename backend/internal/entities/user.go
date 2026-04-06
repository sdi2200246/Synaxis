package entities

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID           uuid.UUID `db:"id"`
    Username     string    `db:"username"`
    PasswordHash string    `db:"password_hash"`
    FirstName    string    `db:"first_name"`
    LastName     string    `db:"last_name"`
    Email        string    `db:"email"`
    Phone        string    `db:"phone"`
    Address      string    `db:"address"`
    City         string    `db:"city"`
    Country      string    `db:"country"`
    TaxID        string    `db:"tax_id"`
    Role         string    `db:"role"`
    Status       string    `db:"status"`
    CreatedAt    time.Time `db:"created_at"`
    UpdatedAt    *time.Time `db:"updated_at"`
}

type UserFilter struct{
	Country  *string
    Status   *string
    CreatedAt *time.Time
}

type UserUpdate struct {
    FirstName *string
    LastName  *string
    Email     *string
    Phone     *string
    Address   *string
    City      *string
    Country   *string
    TaxID     *string
    Role      *string
    Status    *string
}