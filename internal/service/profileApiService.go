package service

import (
	"context"
	"fmt"

	"github.com/eugenshima/trading-api/internal/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type ProfileService struct {
	rps ProfileApiRepository
}

func NewProfileService(rps ProfileApiRepository) *ProfileService {
	return &ProfileService{rps: rps}
}

type ProfileApiRepository interface {
	CreateProfile(context.Context, *model.User) error
}

func (s *ProfileService) Login(context.Context, *model.Auth) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (s *ProfileService) SignUp(ctx context.Context, user *model.User) error {
	hashedPassword := hashPassword([]byte(user.Password))
	user.Password = hashedPassword
	err := s.rps.CreateProfile(ctx, user)
	if err != nil {
		logrus.WithFields(logrus.Fields{"user": user}).Errorf("CreateProfile: %v", err)
		return fmt.Errorf("CreateProfile: %w", err)
	}
	return nil
}

// HashPassword func returns hashed password using bcrypt algorithm
func hashPassword(password []byte) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil
	}
	return hashedPassword
}
