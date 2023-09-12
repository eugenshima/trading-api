// Package repository contains methods to communicate with postgres and gRPC servers
package repository

import (
	"context"
	"fmt"

	"github.com/eugenshima/trading-api/internal/model"

	proto "github.com/eugenshima/profile/proto"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// ProfileRepository represents a repository struct that contains methods to communicate with postgres and gRPCs
type ProfileRepository struct {
	client proto.ProfilesClient
}

// NewProfileRepository creates a new ProfileRepository
func NewProfileRepository(client proto.ProfilesClient) *ProfileRepository {
	return &ProfileRepository{client: client}
}

// Login method to login to the repository server
func (r *ProfileRepository) Login(ctx context.Context, login string, password []byte) (uuid.UUID, error) {
	protoAuth := &proto.Auth{
		Login:    login,
		Password: password,
	}
	response, err := r.client.Login(ctx, &proto.LoginRequest{Auth: protoAuth})
	if err != nil {
		return uuid.Nil, fmt.Errorf("Login: %w", err)
	}
	ID, err := uuid.Parse(response.ID)
	if err != nil {
		logrus.WithFields(logrus.Fields{"ID": response.ID}).Errorf("Parse: %v", err)
		return uuid.Nil, fmt.Errorf("parse: %w", err)
	}
	return ID, nil
}

// UpdateProfile method updates the profile
func (r *ProfileRepository) UpdateProfile(ctx context.Context, id uuid.UUID, refreshToken []byte) error {
	updateToken := &proto.UpdateProfileRequest{
		ID:           id.String(),
		RefreshToken: refreshToken,
	}
	_, err := r.client.UpdateProfile(ctx, updateToken)
	if err != nil {
		return fmt.Errorf("UpdateProfile: %w", err)
	}
	return nil
}

// CreateProfile method creates a profile
func (r *ProfileRepository) CreateProfile(ctx context.Context, profile *model.User) error {
	protoProfile := &proto.CreateProfile{
		Login:    profile.Login,
		Password: profile.Password,
		Username: profile.Username,
	}
	_, err := r.client.CreateNewProfile(ctx, &proto.CreateNewProfileRequest{Profile: protoProfile})
	if err != nil {
		return fmt.Errorf("CreateNewProfile: %w", err)
	}
	return nil
}

// GetRefreshTokenByID method gets a refresh token by given ID
func (r *ProfileRepository) GetRefreshTokenByID(ctx context.Context, ID uuid.UUID) ([]byte, error) {
	response, err := r.client.GetProfileByID(ctx, &proto.GetProfileByIDRequest{ID: ID.String()})
	if err != nil {
		return nil, fmt.Errorf("GetProfileByID: %w", err)
	}
	return response.Profile.RefreshToken, nil
}

// DeleteProfileByID method deletes profile by given ID
func (r *ProfileRepository) DeleteProfileByID(ctx context.Context, ID uuid.UUID) error {
	_, err := r.client.DeleteProfileByID(ctx, &proto.DeleteProfileByIDRequest{ID: ID.String()})
	if err != nil {
		return fmt.Errorf("DeleteProfileByID: %w", err)
	}
	return nil
}
