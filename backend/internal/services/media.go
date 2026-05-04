package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/entities"
	apperr "github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/interfaces"
)

type Media struct {
	ID         uuid.UUID
	EventID    uuid.UUID
	Filename   string
	UploadedAt time.Time
}


type MediaService struct {
	mediaRepo  interfaces.MediaRepository
	eventsRepo interfaces.EventRepository
}

func NewMediaService(mr interfaces.MediaRepository, er interfaces.EventRepository) *MediaService {
	return &MediaService{mediaRepo: mr, eventsRepo: er}
}

func (s *MediaService) Upload(ctx context.Context, callerID, eventID uuid.UUID, size int64, ext string) (Media, error) {
	event, err := s.eventsRepo.GetByID(ctx, eventID)
	if err != nil {
		return Media{}, err
	}
	if err := validateOwnership(callerID, event.OrganizerID); err != nil {
		return Media{}, err
	}

	existing, err := s.mediaRepo.GetByEventID(ctx, eventID)
	if err != nil {
		return Media{}, err
	}
	if len(existing) > 0 {
		return Media{}, fmt.Errorf("event already has a photo: %w", apperr.ErrConflict)
	}

	media := entities.Media{
		ID:         uuid.New(),
		EventID:    eventID,
		Filename:   uuid.New().String() + strings.ToLower(ext),
		SizeBytes:  size,
		UploadedAt: time.Now(),
	}

	if err := media.ApproveCreate() ; err != nil{
		return Media{}, err
	} 

	if err := s.mediaRepo.Create(ctx, media); err != nil {
		return Media{}, err
	}
	return toMedia(media), nil
}

func (s *MediaService) Delete(ctx context.Context, callerID, eventID, mediaID uuid.UUID) (Media, error) {
    event, err := s.eventsRepo.GetByID(ctx, eventID)
    if err != nil {
        return Media{}, err
    }
    if err := validateOwnership(callerID, event.OrganizerID); err != nil {
        return Media{}, err
    }

    photos, err := s.mediaRepo.GetByEventID(ctx, eventID)
    if err != nil {
        return Media{}, err
    }
    var target entities.Media
    for _, p := range photos {
        if p.ID == mediaID {
            target = p
            break
        }
    }
    if target.ID == uuid.Nil {
        return Media{}, apperr.ErrNotFound
    }

    if err := s.mediaRepo.Delete(ctx, mediaID); err != nil {
        return Media{}, err
    }
    return toMedia(target), nil
}


func (s *MediaService) GetByEventID(ctx context.Context, eventID uuid.UUID) ([]entities.Media, error) {
	return s.mediaRepo.GetByEventID(ctx, eventID)
}

func toMedia(m entities.Media) Media {
	return Media{
		ID:         m.ID,
		EventID:    m.EventID,
		Filename:   m.Filename,
		UploadedAt: m.UploadedAt,
	}
}