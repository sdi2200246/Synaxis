package entities

import (
    "time"
    "github.com/google/uuid"
)

type Media struct {
    ID         uuid.UUID `json:"id" db:"id"`
    EventID    uuid.UUID `json:"event_id" db:"event_id"`
    Filename   string    `json:"filename" db:"filename"`
    UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
}
