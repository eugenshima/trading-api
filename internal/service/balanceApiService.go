// Package service contains business-logic methods
package service

import (
	"context"
	"fmt"

	"github.com/eugenshima/trading-api/internal/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// BalanceService struct ....
type BalanceService struct {
	balanceRps BalanceRepository
}

// NewBalanceService creates a new BalanceService
func NewBalanceService(balanceRps BalanceRepository) *BalanceService {
	return &BalanceService{balanceRps: balanceRps}
}

// BalanceRepository interface represents a balance repository
type BalanceRepository interface {
	CreateBalance(context.Context, uuid.UUID) error
	GetBalance(context.Context, uuid.UUID) (*model.Balance, error)
	UpdateBalance(context.Context, *model.Balance) error
}

// GetBalance method gets a balance by given ID
func (s *BalanceService) GetBalance(ctx context.Context, ID uuid.UUID) (*model.Balance, error) {
	return s.balanceRps.GetBalance(ctx, ID)
}

// DepositMoney method adds money to given balance
func (s *BalanceService) DepositMoney(ctx context.Context, balance *model.Balance) (float64, error) {
	dbBalance, err := s.balanceRps.GetBalance(ctx, balance.ProfileID)
	if err != nil {
		return 0, fmt.Errorf("GetBalance: %w", err)
	}
	updatedBalance := addittionSubtractionOperations(dbBalance.Balance, balance.Balance, true)

	dbBalance.Balance = updatedBalance

	err = s.balanceRps.UpdateBalance(ctx, dbBalance)
	if err != nil {
		return 0, fmt.Errorf("UpdateBalance: %w", err)
	}
	return dbBalance.Balance, nil
}

// WithdrawMoney method subs money from given balance
func (s *BalanceService) WithdrawMoney(ctx context.Context, balance *model.Balance) (float64, error) {
	dbBalance, err := s.balanceRps.GetBalance(ctx, balance.ProfileID)
	if err != nil {
		return 0, fmt.Errorf("GetBalance: %w", err)
	}
	updatedBalance := addittionSubtractionOperations(dbBalance.Balance, balance.Balance, false)

	dbBalance.Balance = updatedBalance

	err = s.balanceRps.UpdateBalance(ctx, dbBalance)
	if err != nil {
		return 0, fmt.Errorf("UpdateBalance: %w", err)
	}
	return dbBalance.Balance, nil
}

func (s *BalanceService) CreateBalance(ctx context.Context, profileID uuid.UUID) error {
	err := s.balanceRps.CreateBalance(ctx, profileID)
	if err != nil {
		return fmt.Errorf("CreateBalance: %w", err)
	}
	return nil
}

// addittionSubtractionOperations function calculates balance changes
func addittionSubtractionOperations(dbBalance float64, moneyAmount float64, addSub bool) float64 {
	dbBalanceDecimal := decimal.NewFromFloat(dbBalance)
	moneyAmountDecimal := decimal.NewFromFloat(moneyAmount)

	if addSub {
		dbBalanceDecimal = dbBalanceDecimal.Add(moneyAmountDecimal)
	} else {
		dbBalanceDecimal = dbBalanceDecimal.Sub(moneyAmountDecimal)
	}

	updatedBalance := dbBalanceDecimal.InexactFloat64()

	return updatedBalance
}
