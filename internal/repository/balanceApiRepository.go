package repository

import (
	"context"
	"fmt"

	"github.com/eugenshima/trading-api/internal/model"
	proto "github.com/eugenshima/trading-api/proto/balance"
	"github.com/google/uuid"
)

type BalanceRepository struct {
	client proto.BalanceServiceClient
}

func NewBalanceRepository(client proto.BalanceServiceClient) *BalanceRepository {
	return &BalanceRepository{client: client}
}

func (r *BalanceRepository) GetBalance(ctx context.Context, id uuid.UUID) (*model.Balance, error) {
	response, err := r.client.GetUserByID(ctx, &proto.UserGetByIDRequest{ProfileID: id.String()})
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

func (r *BalanceRepository) UpdateBalance(ctx context.Context, balance *model.Balance) error {
	return nil
}
