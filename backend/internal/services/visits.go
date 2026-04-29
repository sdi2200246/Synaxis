package services

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/entities"
	"github.com/sdi2200246/synaxis/internal/interfaces"
)

type VisitService struct {
	visitRepo 	   interfaces.VisitsRepository
	visitsDetector map[string]time.Time
	mu             sync.RWMutex
}

func NewVisitService(vr interfaces.VisitsRepository) *VisitService {
	return &VisitService{visitRepo: vr , visitsDetector: make(map[string]time.Time)}
}

func (s *VisitService) RecordVisit(ctx context.Context, userID, eventID uuid.UUID) error {
	key := userID.String() + eventID.String()

	s.mu.RLock()
	lastSeen, ok := s.visitsDetector[key]
	s.mu.RUnlock()

	if ok && time.Since(lastSeen) < 10*time.Second {
		return nil
	}

	s.mu.Lock()
	s.visitsDetector[key] = time.Now()
	s.mu.Unlock()

	visit := entities.Visit{
		ID:        uuid.New(),
		UserID:    userID,
		EventID:   eventID,
		VisitedAt: time.Now(),
	}
	return s.visitRepo.Create(ctx, visit)
}