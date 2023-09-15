// Package handlers for handling echo requests
package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	middlewr "github.com/eugenshima/trading-api/internal/middleware"
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
	GetBalance(context.Context, uuid.UUID) (*model.Balance, error)
	DepositMoney(context.Context, *model.Balance) (float64, error)
	WithdrawMoney(context.Context, *model.Balance) (float64, error)
	CreateBalance(context.Context, uuid.UUID) error
}

// Deposit function for adding some amount of money to a balance
func (h *BalanceAPIHandler) Deposit(c echo.Context) error {
	reqBalance := &model.Balance{}
	err := c.Bind(reqBalance)
	if err != nil {
		logrus.WithFields(logrus.Fields{"reqBalance": reqBalance}).Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	id, err := middlewr.GetPayloadFromToken(strings.Split(c.Request().Header.Get("Authorization"), " ")[1])
	if err != nil {
		logrus.WithFields(logrus.Fields{"Payload": strings.Split(c.Request().Header.Get("Authorization"), " ")[1]}).Errorf("GetPayloadFromToken: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("GetPayloadFromToken: %v", err))
	}
	reqBalance.ProfileID = id
	currentBalance, err := h.srv.DepositMoney(c.Request().Context(), reqBalance)
	if err != nil {
		logrus.WithFields(logrus.Fields{"reqBalance": reqBalance}).Errorf("DepositMoney: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("DepositMoney: %v", err))
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("CurrentBalance: %v", currentBalance))
}

// Withdraw function for removing some amount of money from a balance
func (h *BalanceAPIHandler) Withdraw(c echo.Context) error {
	reqBalance := &model.Balance{}
	err := c.Bind(reqBalance)
	if err != nil {
		logrus.WithFields(logrus.Fields{"reqBalance": reqBalance}).Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	id, err := middlewr.GetPayloadFromToken(strings.Split(c.Request().Header.Get("Authorization"), " ")[1])
	if err != nil {
		logrus.WithFields(logrus.Fields{"Payload": strings.Split(c.Request().Header.Get("Authorization"), " ")[1]}).Errorf("GetPayloadFromToken: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("GetPayloadFromToken: %v", err))
	}
	reqBalance.ProfileID = id
	currentBalance, err := h.srv.WithdrawMoney(c.Request().Context(), reqBalance)
	if err != nil {
		logrus.WithFields(logrus.Fields{"reqBalance": reqBalance}).Errorf("WithdrawMoney: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("WithdrawMoney: %v", err))
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("CurrentBalance: %v", currentBalance))
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
	balance, err := h.srv.GetBalance(c.Request().Context(), reqBalance.ProfileID)
	if err != nil {
		logrus.WithFields(logrus.Fields{"ID": reqBalance.ProfileID}).Errorf("GetBalance: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("GetBalance: %v", err))
	}
	return c.JSON(http.StatusOK, balance)
}

func (h *BalanceAPIHandler) CreateBalance(c echo.Context) error {
	id, err := middlewr.GetPayloadFromToken(strings.Split(c.Request().Header.Get("Authorization"), " ")[1])
	if err != nil {
		logrus.WithFields(logrus.Fields{"Payload": strings.Split(c.Request().Header.Get("Authorization"), " ")[1]}).Errorf("GetPayloadFromToken: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("GetPayloadFromToken: %v", err))
	}
	err = h.srv.CreateBalance(c.Request().Context(), id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": id}).Errorf("CreateBalance: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("CreateBalance: %v", err))
	}
	return c.JSON(http.StatusOK, id)
}
