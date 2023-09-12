// Package handlers for handling echo requests
package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eugenshima/trading-api/internal/model"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// BalanceAPIHandler struct represents a handler for Balance API requests
type BalanceAPIHandler struct {
	srv BalanceAPIService
}

// NewBalanceAPIHandler creates a new BalanceApiHandler
func NewBalanceAPIHandler(srv BalanceAPIService) *BalanceAPIHandler {
	return &BalanceAPIHandler{srv: srv}
}

// BalanceAPIService represents a service for Balance API requests
type BalanceAPIService interface {
	AddSubscriber(context.Context, []string) (*model.Shares, error)
	GetBalance(context.Context, uuid.UUID) (*model.Balance, error)
}

// Deposit function for adding some amount of money to a balance
func (h *BalanceAPIHandler) Deposit(c echo.Context) error {
	return c.JSON(http.StatusOK, "Deposit")
}

// Withdraw function for removing some amount of money from a balance
func (h *BalanceAPIHandler) Withdraw(c echo.Context) error {
	return c.JSON(http.StatusOK, "Withdraw")
}

// GetBalance function return balance of given account
// nolint: dupl
func (h *BalanceAPIHandler) GetBalance(c echo.Context) error {
	reqBalance := &model.Balance{}
	err := c.Bind(reqBalance)
	if err != nil {
		logrus.WithFields(logrus.Fields{"reqBalance": reqBalance}).Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	balance, err := h.srv.GetBalance(c.Request().Context(), reqBalance.ID)
	if err != nil {
		logrus.WithFields(logrus.Fields{"ID": reqBalance.ID}).Errorf("GetBalance: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("GetBalance: %v", err))
	}
	return c.JSON(http.StatusOK, balance)
}

// GetLatestPrice function return latest price for chosen share
// nolint: dupl
func (h *BalanceAPIHandler) GetLatestPrice(c echo.Context) error {
	streamedShares := &model.StreamedShares{}
	err := c.Bind(streamedShares)
	if err != nil {
		logrus.WithFields(logrus.Fields{"streamedShares": streamedShares}).Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	share, err := h.srv.AddSubscriber(c.Request().Context(), streamedShares.Share)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Shares": streamedShares.Share}).Errorf("AddSubscriber: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("AddSubscriber: %v", err))
	}
	return c.JSON(http.StatusOK, share)
}
