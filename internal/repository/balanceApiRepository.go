// Package repository contains methods to communicate with postgres and gRPC servers
package repository

import (
	"context"
	"fmt"

	balanceProto "github.com/eugenshima/balance/proto"
	"github.com/eugenshima/trading-api/internal/model"
	"github.com/google/uuid"
)

// BalanceRepository struct represents a repository
type BalanceRepository struct {
	client balanceProto.BalanceServiceClient
}

// NewBalanceRepository creates a new BalanceRepository
func NewBalanceRepository(client balanceProto.BalanceServiceClient) *BalanceRepository {
	return &BalanceRepository{client: client}
}

func (r *BalanceRepository) CreateBalance(ctx context.Context, profileID uuid.UUID) error {
	protoBalance := &balanceProto.Balance{
		ProfileID: profileID.String(),
		Balance:   0,
	}
	_, err := r.client.CreateUserBalance(ctx, &balanceProto.CreateBalanceRequest{Balance: protoBalance})
	if err != nil {
		return fmt.Errorf("CreateUserBalance: %w", err)
	}
	return nil
}

// GetBalance method returns a balance by the given ID
func (r *BalanceRepository) GetBalance(ctx context.Context, id uuid.UUID) (*model.Balance, error) {
	response, err := r.client.GetUserByID(ctx, &balanceProto.UserGetByIDRequest{ProfileID: id.String()})
	if err != nil {
		return nil, fmt.Errorf("GetUserByID: %w", err)
	}
	responseProfileID, err := uuid.Parse(response.Balance.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	responseBalanceID, err := uuid.Parse(response.Balance.BalanceID)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	balance := &model.Balance{
		BalanceID: responseBalanceID,
		ProfileID: responseProfileID,
		Balance:   response.Balance.Balance,
	}
	return balance, nil
}

// UpdateBalance method updates a balance
func (r *BalanceRepository) UpdateBalance(ctx context.Context, balance *model.Balance) error {
	protoBalance := &balanceProto.Balance{
		BalanceID: balance.BalanceID.String(),
		ProfileID: balance.ProfileID.String(),
		Balance:   balance.Balance,
	}
	_, err := r.client.UpdateUserBalance(ctx, &balanceProto.UserUpdateRequest{Balance: protoBalance})
	if err != nil {
		return fmt.Errorf("UpdateUserBalance: %w", err)
	}
	return nil
}
