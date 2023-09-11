package service

import (
	"context"

	"github.com/eugenshima/trading-api/internal/model"
	"github.com/google/uuid"
)

type BalanceService struct {
	balanceRps      BalanceRepository
	priceServiceRps PriceServiceRepository
}

func NewBalanceService(balanceRps BalanceRepository, priceServiceRps PriceServiceRepository) *BalanceService {
	return &BalanceService{balanceRps: balanceRps, priceServiceRps: priceServiceRps}
}

type BalanceRepository interface {
	GetBalance(context.Context, uuid.UUID) (*model.Balance, error)
	UpdateBalance(context.Context, *model.Balance) error
}

type PriceServiceRepository interface {
	RecvShares(context.Context, []string) (*model.Shares, error)
}

func (s *BalanceService) AddSubscriber(ctx context.Context, shares []string) (*model.Shares, error) {
	return s.priceServiceRps.RecvShares(ctx, shares)
}

func (s *BalanceService) GetBalance(ctx context.Context, ID uuid.UUID) (*model.Balance, error) {
	return s.balanceRps.GetBalance(ctx, ID)
}
