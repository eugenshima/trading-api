package repository

import (
	"context"
	"fmt"

	"github.com/eugenshima/trading-api/internal/model"
	proto "github.com/eugenshima/trading-api/proto"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ProfileRepository struct {
	client proto.PriceServiceClient
}

func NewProfileRepository(client proto.PriceServiceClient) *ProfileRepository {
	return &ProfileRepository{client: client}
}

func (r *ProfileRepository) Login(ctx context.Context, login string, password []byte) (uuid.UUID, error) {
	protoAuth := &proto.Auth{
		Login:    login,
		Password: password,
	}
	response, err := r.client.Login(ctx, &proto.LoginRequest{Auth: protoAuth})
	if err != nil {
		logrus.WithFields(logrus.Fields{"protoAuth": protoAuth}).Errorf("Login: %v", err)
		return uuid.Nil, fmt.Errorf("Login: %w", err)
	}
	ID, err := uuid.Parse(response.ID)
	if err != nil {
		logrus.WithFields(logrus.Fields{"ID": response.ID}).Errorf("Parse: %v", err)
		return uuid.Nil, fmt.Errorf("parse: %w", err)
	}
	return ID, nil
}

func (r *ProfileRepository) CreateProfile(ctx context.Context, profile *model.User) error {
	protoProfile := &proto.CreateProfile{
		Login:    profile.Login,
		Password: profile.Password,
		Username: profile.Username,
	}
	_, err := r.client.CreateNewProfile(ctx, &proto.CreateNewProfileRequest{Profile: protoProfile})
	if err != nil {
		logrus.WithFields(logrus.Fields{"protoProfile": protoProfile}).Errorf("CreateNewProfile: %v", err)
		return fmt.Errorf("CreateNewProfile: %w", err)
	}
	return nil
}
