// Package service contains business-logic methods
package service

import (
	"context"

	"github.com/eugenshima/trading-api/internal/model"
	"github.com/google/uuid"
)

// BalanceService struct ....
type BalanceService struct {
	balanceRps      BalanceRepository
	priceServiceRps PriceServiceRepository
}

// NewBalanceService creates a new BalanceService
func NewBalanceService(balanceRps BalanceRepository, priceServiceRps PriceServiceRepository) *BalanceService {
	return &BalanceService{balanceRps: balanceRps, priceServiceRps: priceServiceRps}
}

// BalanceRepository interface represents a balance repository
type BalanceRepository interface {
	GetBalance(context.Context, uuid.UUID) (*model.Balance, error)
	UpdateBalance(context.Context, *model.Balance) error
}

// PriceServiceRepository interface represents a price service repository
type PriceServiceRepository interface {
	RecvShares(context.Context, []string) (*model.Shares, error)
}

// AddSubscriber method adds a new subscriber to price service
func (s *BalanceService) AddSubscriber(ctx context.Context, shares []string) (*model.Shares, error) {
	return s.priceServiceRps.RecvShares(ctx, shares)
}

// GetBalance method gets a balance by given ID
func (s *BalanceService) GetBalance(ctx context.Context, ID uuid.UUID) (*model.Balance, error) {
	return s.balanceRps.GetBalance(ctx, ID)
}
