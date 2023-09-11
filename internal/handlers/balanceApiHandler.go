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

type BalanceApiHandler struct {
	srv BalanceApiService
}

func NewBalanceApiHandler(srv BalanceApiService) *BalanceApiHandler {
	return &BalanceApiHandler{srv: srv}
}

type BalanceApiService interface {
	AddSubscriber(context.Context, []string) (*model.Shares, error)
	GetBalance(context.Context, uuid.UUID) (*model.Balance, error)
}

func (h *BalanceApiHandler) Deposit(c echo.Context) error {
	return c.JSON(http.StatusOK, "Deposit")
}

func (h *BalanceApiHandler) Withdraw(c echo.Context) error {
	return c.JSON(http.StatusOK, "Withdraw")
}

// GetBalance function return balance of given account
func (h *BalanceApiHandler) GetBalance(c echo.Context) error {
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
func (h *BalanceApiHandler) GetLatestPrice(c echo.Context) error {
	streamedShares := &model.StreamedShares{}
	err := c.Bind(streamedShares)
	if err != nil {
		logrus.WithFields(logrus.Fields{"streamedShares": streamedShares}).Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	share, err := h.srv.AddSubscriber(c.Request().Context(), streamedShares.Shares)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Shares": streamedShares.Shares}).Errorf("AddSubscriber: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("AddSubscriber: %v", err))
	}
	return c.JSON(http.StatusOK, share)
}
