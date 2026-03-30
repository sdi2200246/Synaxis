package services

import (
	"context"
	"testing"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/entities"
	apperr "github.com/sdi2200246/synaxis/internal/error"
)

// MockUserRepo implements interfaces.UserRepository
type MockUserRepo struct {
	users []entities.User
}

func (m *MockUserRepo) Create(ctx context.Context, user entities.User) error {
	return nil // stub
}

func (m *MockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (entities.User, error) {
	return entities.User{}, nil // stub
}

func (m *MockUserRepo) GetByUsername(ctx context.Context, username string) (entities.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return entities.User{}, apperr.ErrNotFound
}

func TestLogin_Success(t *testing.T) {
	// pre-hash the password the same way the service does
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)

	mock := &MockUserRepo{
		users: []entities.User{
			{
				ID:           uuid.New(),
				Username:     "jason",
				PasswordHash: string(hash),
				Role:         "user",
				Status:       "approved",
			},
		},
	}

	svc := NewAuthService(mock, "testsecret")

	token, err := svc.Login(context.Background(), UserCridentials{
		Username: "jason",
		Password: "secret123",
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if token == "" {
		t.Error("expected token, got empty string")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)

	mock := &MockUserRepo{
		users: []entities.User{
			{
				ID:           uuid.New(),
				Username:     "jason",
				PasswordHash: string(hash),
				Role:         "user",
				Status:       "approved",
			},
		},
	}

	svc := NewAuthService(mock, "testsecret")

	_, err := svc.Login(context.Background(), UserCridentials{
		Username: "jason",
		Password: "wrongpassword",
	})

	if err != apperr.ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	mock := &MockUserRepo{}
	svc := NewAuthService(mock, "testsecret")

	_, err := svc.Login(context.Background(), UserCridentials{
		Username: "nobody",
		Password: "secret123",
	})

	if err != apperr.ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestLogin_PendingApproval(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)

	mock := &MockUserRepo{
		users: []entities.User{
			{
				ID:           uuid.New(),
				Username:     "jason",
				PasswordHash: string(hash),
				Role:         "user",
				Status:       "pending",
			},
		},
	}

	svc := NewAuthService(mock, "testsecret")

	_, err := svc.Login(context.Background(), UserCridentials{
		Username: "jason",
		Password: "secret123",
	})

	if err != apperr.ErrPendingApproval {
		t.Errorf("expected ErrPendingApproval, got %v", err)
	}
}
