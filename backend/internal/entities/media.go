package entities

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	apperr "github.com/sdi2200246/synaxis/internal/error"
)

var allowedExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
}
const (
	maxPhotoBytes  = 5 << 20 // 5 MiB
)

type Media struct {
    ID         uuid.UUID `db:"id"`
    EventID    uuid.UUID `db:"event_id"`
    Filename   string    `db:"filename"`
    SizeBytes  int64     `db:"size_bytes"`
    UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
}


func (m Media) ApproveCreate() error {
    if m.SizeBytes > maxPhotoBytes {
        return fmt.Errorf("photo exceeds %d bytes: %w", maxPhotoBytes, apperr.ErrBadInput)
    }
    ext := strings.ToLower(filepath.Ext(m.Filename))
    if !allowedExtensions[ext] {
        return fmt.Errorf("unsupported file extension %q: %w", ext, apperr.ErrBadInput)
    }
    return nil
}