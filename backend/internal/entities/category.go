package entities

import "github.com/google/uuid"

type Category struct {
    ID       uuid.UUID  `json:"id" db:"id"`
    Name     string     `json:"name" db:"name"`
    ParentID *uuid.UUID `json:"parent_id" db:"parent_id"`
}
