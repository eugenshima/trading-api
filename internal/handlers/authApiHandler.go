// Package handlers for the various types of events
package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	middlewr "github.com/eugenshima/trading-api/internal/middleware"
	"github.com/eugenshima/trading-api/internal/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

// ProfileAPIHandler is the handler for the various types of events
type ProfileAPIHandler struct {
	srv ProfileAPIService
}

// NewProfileAPIHandler creates a new ProfileApiHandler
func NewProfileAPIHandler(srv ProfileAPIService) *ProfileAPIHandler {
	return &ProfileAPIHandler{srv: srv}
}

// ProfileAPIService represents profile-api-service
type ProfileAPIService interface {
	Login(context.Context, *model.Login) (*model.JWTResponse, error)
	SignUp(context.Context, *model.User) error
	RefreshTokenPair(context.Context, uuid.UUID, string, []byte) (*model.JWTResponse, error)
	DeleteProfile(context.Context, uuid.UUID) error
}

// Login method handles authentification process
func (h *ProfileAPIHandler) Login(c echo.Context) error {
	login := &model.Login{}
	err := c.Bind(login)
	if err != nil {
		logrus.WithFields(logrus.Fields{"login": login}).Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	jwtResponse, err := h.srv.Login(c.Request().Context(), login)
	if err != nil {
		logrus.WithFields(logrus.Fields{"login": login}).Errorf("Login: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Login: %v", err))
	}
	return c.JSON(http.StatusOK, jwtResponse)
}

// SignUp method to sign up the user to system
func (h *ProfileAPIHandler) SignUp(c echo.Context) error {
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
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("SignUp: %v", err))
	}
	return c.JSON(http.StatusOK, "created")
}

// RefreshTokenPair method refreshes token pair
func (h *ProfileAPIHandler) RefreshTokenPair(c echo.Context) error {
	accessRefresh := &model.JWTResponse{}
	err := c.Bind(accessRefresh)
	if err != nil {
		logrus.WithFields(logrus.Fields{"accessRefresh": accessRefresh}).Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	refreshedTokens, err := h.srv.RefreshTokenPair(c.Request().Context(), accessRefresh.ID, accessRefresh.AccessToken, accessRefresh.RefreshToken)
	if err != nil {
		logrus.WithFields(logrus.Fields{"ID": accessRefresh.ID, "AccessToken": accessRefresh.AccessToken, "RefreshToken": accessRefresh.RefreshToken}).Errorf("RefreshTokenPair: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("RefreshTokenPair: %v", err))
	}
	return c.JSON(http.StatusOK, refreshedTokens)
}

// DeleteProfile method deletes the profile associated with the given ID from token payload
func (h *ProfileAPIHandler) DeleteProfile(c echo.Context) error {
	id, err := middlewr.GetPayloadFromToken(strings.Split(c.Request().Header.Get("Authorization"), " ")[1])
	if err != nil {
		logrus.WithFields(logrus.Fields{"Payload": strings.Split(c.Request().Header.Get("Authorization"), " ")[1]}).Errorf("GetPayloadFromToken: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("GetPayloadFromToken: %v", err))
	}
	err = h.srv.DeleteProfile(c.Request().Context(), id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"ID": id}).Errorf("DeleteProfile: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("DeleteProfile: %v", err))
	}
	return c.JSON(http.StatusOK, "deleted")
}
