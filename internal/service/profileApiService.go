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
	Login(context.Context, string, []byte) (uuid.UUID, error)
}

func (s *ProfileService) Login(ctx context.Context, auth *model.Auth) (uuid.UUID, error) {
	hashedPassword := hashPassword([]byte(auth.Password))
	id, err := s.rps.Login(ctx, auth.Login, hashedPassword)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Login": auth.Login, "hashedPassword": hashedPassword}).Errorf("Login: %v", err)
		return uuid.Nil, fmt.Errorf("Login: %w", err)
	}
	return id, nil
}

func (s *ProfileService) SignUp(ctx context.Context, user *model.User) error {
	hashedPassword := hashPassword(user.Password)
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
