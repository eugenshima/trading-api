package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eugenshima/trading-api/internal/model"
	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

type ProfileApiHandler struct {
	srv ProfileApiService
}

func NewProfileApiHandler(srv ProfileApiService) *ProfileApiHandler {
	return &ProfileApiHandler{srv: srv}
}

type ProfileApiService interface {
	Login(context.Context, *model.Login) (*model.JWTResponse, error)
	SignUp(context.Context, *model.User) error
}

func (h *ProfileApiHandler) Login(c echo.Context) error {
	login := &model.Login{}
	err := c.Bind(login)
	if err != nil {
		logrus.WithFields(logrus.Fields{"login": login}).Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	jwtResponse, err := h.srv.Login(c.Request().Context(), login)
	if err != nil {
		logrus.Info("wrong password")
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	return c.JSON(http.StatusOK, jwtResponse)
}

func (h *ProfileApiHandler) SignUp(c echo.Context) error {
	newUser := &model.NewUser{}
	err := c.Bind(newUser)
	if err != nil {
		logrus.WithFields(logrus.Fields{"user": newUser}).Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	user := &model.User{
		Login:    newUser.Login,
		Password: []byte(newUser.Password),
		Username: newUser.Username,
	}
	err = h.srv.SignUp(c.Request().Context(), user)
	if err != nil {
		logrus.Info("wrong password")
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	return c.JSON(http.StatusOK, "created")
}
