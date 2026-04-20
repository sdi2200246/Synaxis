package services

import (
	"context"
	"time"
	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/entities"
	"github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/interfaces"
	"golang.org/x/crypto/bcrypt"
)

type CandidateUser struct {
    Username  string
    Password  string
    FirstName string
    LastName  string
    Email     string
    Phone     string
    Address   string
    City      string
    Country   string
    TaxID     string
}

type User struct {
    ID        uuid.UUID 
    Username  string    
    FirstName string    
    LastName  string    
    Email     string
    Address    string    
    City      string    
    Country   string
    TaxID     string
    Status    string
    Phone     string    
    CreatedAt time.Time
    UpdatedAt *time.Time
}

type PublicUser struct {
	ID        uuid.UUID
	Username  string
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

type UserFilter struct {
    Country  *string
    Status   *string
    CreatedAt *time.Time
}


type UserService struct{
	userRepo interfaces.UserRepository
}

func NewUserService(r interfaces.UserRepository)*UserService{
	return  &UserService{userRepo: r}
}

func (s UserService)RegisterUser(ctx context.Context , user CandidateUser)error{

    passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return apperr.ErrInternal
    }

	newUser := entities.User{
        ID:           uuid.New(),
        Username:     user.Username,
        PasswordHash: string(passwordHash),
        FirstName:    user.FirstName,
        LastName:     user.LastName,
        Email:        user.Email,
        Phone:        user.Phone,
        Address:      user.Address,
        City:         user.City,
        Country:      user.Country,
        TaxID:        user.TaxID,
        Role:         "user",
        Status:       "pending",
        CreatedAt:    time.Now(),
    }	

	return s.userRepo.Create(ctx, newUser)
}

func (s *UserService) GetUsers(ctx context.Context, f UserFilter) ([]User, error) {
    filter := entities.UserFilter{
        Country:   f.Country,
        Status:    f.Status,
        CreatedAt: f.CreatedAt,
    }

    users, err := s.userRepo.ListUsers(ctx, filter)
    if err != nil {
        return nil, err
    }

    plain := make([]User, len(users))
    for i, u := range users {
        plain[i] = User{
            ID:        u.ID,
            Username:  u.Username,
            FirstName: u.FirstName,
            LastName:  u.LastName,
            Email:     u.Email,
            Phone:     u.Phone,
            Address:   u.Address,
            City:      u.City,
            Country:   u.Country,
            TaxID:     u.TaxID,
            Status:    u.Status,
            CreatedAt: u.CreatedAt,
            UpdatedAt: u.UpdatedAt,
        }
    }

    return plain, nil
}

func (s *UserService) GetPublicByID(ctx context.Context, id uuid.UUID) (PublicUser, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return PublicUser{}, err
	}
	return PublicUser{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
	}, nil
}


func (s *UserService) ApproveUser(ctx context.Context, id uuid.UUID) error {
    status := "approved"
    return s.userRepo.UpdateUser(ctx, id, entities.UserUpdate{
        Status: &status,
    })
}

func (s *UserService) RejectUser(ctx context.Context, id uuid.UUID) error {
    status := "rejected"
    return s.userRepo.UpdateUser(ctx, id, entities.UserUpdate{
        Status: &status,
    })
}