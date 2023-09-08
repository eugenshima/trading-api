package repository

import (
	"context"
	"fmt"

	"github.com/eugenshima/trading-api/internal/model"
	proto "github.com/eugenshima/trading-api/proto"
	"github.com/sirupsen/logrus"
)

type ProfileRepository struct {
	client proto.PriceServiceClient
}

func NewProfileRepository(client proto.PriceServiceClient) *ProfileRepository {
	return &ProfileRepository{client: client}
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
