package services

import (
	"context"
	"time"
	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/entities"
	"github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/interfaces"
	"github.com/sdi2200246/synaxis/internal/types"
	"golang.org/x/crypto/bcrypt"
)

type CandidateUser struct {
    Username  string `json:"username"   binding:"required"`
    Password  string `json:"password"   binding:"required"`
    FirstName string `json:"first_name" binding:"required"`
    LastName  string `json:"last_name"  binding:"required"`
    Email     string `json:"email"      binding:"required,email"`
    Phone     string `json:"phone"      binding:"required"`
    Address   string `json:"address"    binding:"required"`
    City      string `json:"city"       binding:"required"`
    Country   string `json:"country"    binding:"required"`
    TaxID     string `json:"tax_id"     binding:"required"`
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

func (s UserService)GetUsers(ctx context.Context , f types.UserFilter)([]entities.User, error){
	return s.userRepo.ListUsers(ctx , f)
}
