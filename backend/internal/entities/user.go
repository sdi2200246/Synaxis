package entities

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID           uuid.UUID `json:"id"            db:"id"`
    Username     string    `json:"username"       db:"username"`
    PasswordHash string    `json:"-"              db:"password_hash"`
    FirstName    string    `json:"first_name"     db:"first_name"`
    LastName     string    `json:"last_name"      db:"last_name"`
    Email        string    `json:"email"          db:"email"`
    Phone        string    `json:"phone"          db:"phone"`
    Address      string    `json:"address"        db:"address"`
    City         string    `json:"city"           db:"city"`
    Country      string    `json:"country"        db:"country"`
    TaxID        string    `json:"tax_id"         db:"tax_id"`
    Role         string    `json:"role"           db:"role"`
    Status       string    `json:"status"         db:"status"`
    CreatedAt    time.Time `json:"created_at"     db:"created_at"`
}


