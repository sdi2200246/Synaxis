package repos

import (
	"context"
	"testing"

	"github.com/sdi2200246/synaxis/internal/entities"
	"github.com/sdi2200246/synaxis/internal/services"
	apperr "github.com/sdi2200246/synaxis/internal/error"
	"github.com/google/uuid"
	"time"
)

// MockUserRepo implements interfaces.UserRepository
type MockUserRepo struct {
	users []entities.User
}

func (m *MockUserRepo) Create(ctx context.Context, user entities.User) error {
	for _, u := range m.users {
		if u.Username == user.Username {
			return apperr.ErrUsernameConflict
		}
		if u.Email == user.Email {
			return apperr.ErrEmailConflict
		}
	}
	m.users = append(m.users, user)
	return nil
}

func (m *MockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (entities.User, error) {
    return entities.User{}, nil  // stub — not tested yet
}

func (m *MockUserRepo) GetByUsername(ctx context.Context, username string) (entities.User, error) {
    return entities.User{}, nil  // stub — not tested yet
}

func TestRegisterUser_Success(t *testing.T) {
	mock := &MockUserRepo{}
	svc := services.NewUserService(mock)

	err := svc.RegisterUser(context.Background(), services.CandidateUser{
		Username:  "jason",
		Password:  "secret",
		FirstName: "Jason",
		LastName:  "Test",
		Email:     "jason@test.com",
		Phone:     "123456789",
		Address:   "Test St 1",
		City:      "Athens",
		Country:   "Greece",
		TaxID:     "123456789",
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(mock.users) != 1 {
		t.Errorf("expected 1 user in mock, got %d", len(mock.users))
	}

	if mock.users[0].Username != "jason" {
		t.Errorf("expected username jason, got %s", mock.users[0].Username)
	}

	if mock.users[0].Role != "user" {
		t.Errorf("expected role user, got %s", mock.users[0].Role)
	}

	if mock.users[0].Status != "pending" {
		t.Errorf("expected status pending, got %s", mock.users[0].Status)
	}
}

func TestRegisterUser_UsernameConflict(t *testing.T) {
	mock := &MockUserRepo{
		users: []entities.User{
			{
				ID:       uuid.New(),
				Username: "jason",
				Email:    "other@test.com",
				Role:     "user",
				Status:   "pending",
				CreatedAt: time.Now(),
			},
		},
	}
	svc := services.NewUserService(mock)

	err := svc.RegisterUser(context.Background(), services.CandidateUser{
		Username: "jason",
		Email:    "jason@test.com",
	})

	if err != apperr.ErrUsernameConflict {
		t.Errorf("expected ErrUsernameConflict, got %v", err)
	}
}

func TestRegisterUser_EmailConflict(t *testing.T) {
	mock := &MockUserRepo{
		users: []entities.User{
			{
				ID:       uuid.New(),
				Username: "other",
				Email:    "jason@test.com",
				Role:     "user",
				Status:   "pending",
				CreatedAt: time.Now(),
			},
		},
	}
	svc := services.NewUserService(mock)

	err := svc.RegisterUser(context.Background(), services.CandidateUser{
		Username: "jason",
		Email:    "jason@test.com",
	})

	if err != apperr.ErrEmailConflict {
		t.Errorf("expected ErrEmailConflict, got %v", err)
	}
}