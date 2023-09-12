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

// GetBalance method returns a balance by the given ID
func (r *BalanceRepository) GetBalance(ctx context.Context, id uuid.UUID) (*model.Balance, error) {
	response, err := r.client.GetUserByID(ctx, &balanceProto.UserGetByIDRequest{ProfileID: id.String()})
	if err != nil {
		return nil, fmt.Errorf("GetUserByID: %w", err)
	}
	responseID, err := uuid.Parse(response.Balance.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	balance := &model.Balance{
		ID:      responseID,
		Balance: response.Balance.Balance,
	}
	return balance, nil
}

// UpdateBalance method updates a balance
func (r *BalanceRepository) UpdateBalance(_ context.Context, _ *model.Balance) error {
	return nil
}
