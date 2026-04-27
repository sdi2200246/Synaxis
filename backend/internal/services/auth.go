package services

import (
    "context"
    "time"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "github.com/sdi2200246/synaxis/internal/interfaces"
    apperr "github.com/sdi2200246/synaxis/internal/error"
)

type UserCridentials struct{
	Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type Claims struct {
    UserID uuid.UUID `json:"user_id"`
    Role   string   `json:"role"`
    jwt.RegisteredClaims
}

type AuthService struct{
	userRepo interfaces.UserRepository
	secret string
}
func NewAuthService(userRepo interfaces.UserRepository, secret string) *AuthService {
    return &AuthService{userRepo: userRepo, secret: secret}
}

func (s *AuthService)Login(ctx context.Context , credentials UserCridentials)(string , error){
	user , err :=  s.userRepo.GetByUsername(ctx , credentials.Username)

	if err != nil {
		return "" ,apperr.ErrUnauthorized	
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password))
	if err != nil {
    	return "", apperr.ErrUnauthorized
	}

	if user.Status == "pending" {
		return "", apperr.ErrPendingApproval
	}
	if user.Status == "rejected" {
		return "", apperr.ErrRejected
	}

    claims := Claims{
        UserID: user.ID,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(s.secret))
    if err != nil {
        return "", apperr.ErrInternal
    }

    return tokenString, nil

}

func validateOwnership(callerID, ownerID uuid.UUID) error {
    if callerID != ownerID {
        return apperr.ErrForbidden
    }
    return nil
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
        return []byte(s.secret), nil
    })
    if err != nil || !token.Valid {
        return nil, apperr.ErrUnauthorized
    }
    claims, ok := token.Claims.(*Claims)
    if !ok {
        return nil, apperr.ErrUnauthorized
    }
    return claims, nil
}